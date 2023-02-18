package logic

import (
	"douyin5856/dao/mysql"
	"douyin5856/dao/redis"
	"douyin5856/models"
	"go.uber.org/zap"
	"log"
	"sort"
	"sync"
)

// FavoriteAction 处理点赞相关逻辑
func FavoriteAction(videoId, userId int64, actionType int32) (err error) {
	// 1、查找数据库有没有点赞或者不点赞的记录
	haven, isFavorite, err := mysql.QueryFavorite(videoId, userId)
	if err != nil {
		return err
	}
	// 2、如果点赞过，直接bool取反，修改数据库即可
	if haven {
		isFavorite = !isFavorite
		if err = mysql.AlterFavorite(userId, videoId, isFavorite); err != nil {
			return err
		}
		if actionType == 1 {
			zap.L().Info("点赞点赞")
			if err = mysql.AddFavoriteCount(videoId, 1); err != nil {
				return err
			}
		} else {
			zap.L().Info("取消点赞")
			if err = mysql.AddFavoriteCount(videoId, -1); err != nil {
				return err
			}
		}
	} else {
		// 3、如果没有查到，则插入新数据
		if actionType == 1 {
			zap.L().Info("点赞点赞")
			isFavorite = true
			if err = mysql.AddFavoriteCount(videoId, 1); err != nil {
				return err
			}
		}
		if err = mysql.InsertFavorite(userId, videoId, isFavorite); err != nil {
			return err
		}
	}

	// 上面是对数据库的操作，下面我们需要更新redis中的数据
	if actionType == 1 {
		//增加用户的喜欢集合
		redis.AddUserFavoriteSet(userId, videoId)
	} else {
		//把视频从用户喜欢集合里面删除
		redis.DelUserFavoriteSet(userId, videoId)
	}

	return nil
}

// FavoriteList 处理点赞列表相关逻辑
func FavoriteList(userId, curUserId int64) (*models.ResponseFavoriteList, error) {
	log.Println("FavoriteList : running")
	// 1、更具userid从redis中读取喜欢的视频
	videoIdList, err := redis.QueryFavoriteSet(userId)
	if err != nil {
		return nil, err
	}
	if len(videoIdList) == 0 {
		// 1、根据userId从user_favorite_video表中找到点赞的视频id
		videoIdList, err = mysql.QueryFavoriteVideos(userId)
		if err != nil {
			return nil, err
		}
	}
	// 2、通过视频id 批量查找视频数据
	videosTables, err := mysql.QueryVideos(videoIdList)
	if err != nil {
		return nil, err
	}
	// 3、批量查找每个视频作者相关信息
	responseVideos := make([]models.Video, len(videoIdList))
	wgFavoriteList := &sync.WaitGroup{}
	wgFavoriteList.Add(len(videoIdList))
	num := 0
	for _, videosTable := range videosTables {
		var responseVideo models.Video
		go func(videoTable models.VideosTable) {
			// 填充一条视频数据
			StuffOneVideo(&responseVideo, &videoTable, curUserId)
			responseVideos[num] = responseVideo
			num = num + 1
			wgFavoriteList.Done()
		}(videosTable)
	}
	wgFavoriteList.Wait()
	// 因为没有点赞时间这个参数，所以我们按照视频id进行排序
	sort.Sort(FeedSlice(responseVideos))
	// 6、拼接数据
	responseFavoriteList := &models.ResponseFavoriteList{
		models.Response{
			0,
			"请求视频成功",
		},
		responseVideos,
	}
	return responseFavoriteList, nil
}
