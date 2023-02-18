package redis

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"strconv"
)

// QueryFavoriteRdb redis查询用户是否喜欢该视频
func QueryFavoriteRdb(userId, videoId int64) (bool, error) {
	userIds := strconv.FormatInt(userId, 10)
	videoIds := strconv.FormatInt(videoId, 10)
	// 查询redis中该用户是否存在
	if result, err := rdbFavorite.Exists(userIds).Result(); result > 0 {
		if err != nil {
			zap.L().Error("redis/user QueryFavoriteRdb rdbFavorite.Exists failed", zap.Error(err))
			return false, err
		}

		// 判断该用户是否喜欢该视频
		favoriteOr, err1 := rdbFavorite.SIsMember(userIds, videoIds).Result()
		if err1 != nil {
			zap.L().Error("redis/user QueryFavoriteRdb rdbFavorite.SIsMember failed", zap.Error(err))
			return false, err1
		}
		//log.Println("QueryFavoriteRdb success")
		return favoriteOr, nil
	}
	// 用户不存在
	// redis 中该用户不存在，我们添加该用户，并且添加一个默认值
	if _, err := rdbFavorite.SAdd(userIds, -1).Result(); err != nil {
		zap.L().Error("redis/user QueryFavoriteRdb rdbFavorite.SAdd failed", zap.Error(err))
		rdbFavorite.Del(userIds)
		return false, err
	}
	// 设置该用户的一个有效期
	if _, err := rdbFavorite.Expire(userIds, viper.GetDuration("redis.expire_time_long")).Result(); err != nil {
		zap.L().Error("redis/user QueryFavoriteRdb rdbFavorite.Expire failed", zap.Error(err))
		rdbFavorite.Del(userIds)
		return false, err
	}
	return false, nil
}

// UpdateFavoriteRdb 更新redis中用户喜欢视频列表
func UpdateFavoriteRdb(userId int64, videoIds []int64) {
	// 笨蛋版的循环更新
	for _, videoId := range videoIds {
		rdbFavorite.SAdd(strconv.FormatInt(userId, 10), strconv.FormatInt(videoId, 10))
	}
}

// AddUserFavoriteSet 把视频id增加到 用户喜欢集合
func AddUserFavoriteSet(userId, videoId int64) {
	rdbFavorite.SAdd(strconv.FormatInt(userId, 10), strconv.FormatInt(videoId, 10))
	log.Println("视频加入redis中点赞列表")
}

// DelUserFavoriteSet 把视频id从用户喜欢集合删除
func DelUserFavoriteSet(userId, videoId int64) {
	rdbFavorite.SRem(strconv.FormatInt(userId, 10), strconv.FormatInt(videoId, 10))
	log.Println("从redis中删除用户取消点赞的视频")
}

// QueryFavoriteSet 根据用户id从redis中查找喜欢的视频列表
func QueryFavoriteSet(userId int64) ([]int64, error) {
	strings, err := rdbFavorite.SMembers(strconv.FormatInt(userId, 10)).Result()
	res := make([]int64, len(strings))
	for i, str := range strings {
		res[i], _ = strconv.ParseInt(str, 10, 64)
	}
	return res, err
}
