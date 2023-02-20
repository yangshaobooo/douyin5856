package mysql

import (
	"database/sql"
	"douyin5856/models"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"strings"
)

// QueryFavorite 查看当前用户是否有该视频的喜欢记录
func QueryFavorite(videoId, userId int64) (bool, bool, error) {
	sqlStr := `select is_favorite 
	from user_favorite_video 
	use index(idx_user_video)
	where user_id=? and video_id=?`
	var isFavorite bool
	if err := db.Get(&isFavorite, sqlStr, userId, videoId); err != nil {
		if err == sql.ErrNoRows {
			//没有查到就是没有过记录
			return false, isFavorite, nil
		} else {
			zap.L().Error("mysql/QueryFavorite Get failed", zap.Error(err))
			return false, isFavorite, err
		}
	}
	return true, isFavorite, nil
}

// AlterFavorite 修改数据库中的点赞状态
func AlterFavorite(userId, videoId int64, isFavorite bool) (err error) {
	sqlStr := `update user_favorite_video 
	set is_favorite=? 
	where user_id=? and video_id=?`
	_, err = db.Exec(sqlStr, isFavorite, userId, videoId)
	if err != nil {
		zap.L().Error("mysql/AlterFavorite Exec failed", zap.Error(err))
		return err
	}
	//修改点赞状态成功
	return nil
}

// InsertFavorite 向数据库中插入用户是否喜欢的视频
func InsertFavorite(userId, videoId int64, isFavorite bool) (err error) {
	sqlStr := `insert into user_favorite_video
	(user_id,video_id,is_favorite)
	values(?,?,?)`
	_, err = db.Exec(sqlStr, userId, videoId, isFavorite)
	if err != nil {
		zap.L().Error("mysql/InsertFavorite Exec failed", zap.Error(err))
		return err
	}
	return nil
}

// AddVideoFavoriteCount 视频中的点赞数+1
func AddVideoFavoriteCount(videoId, num int64) (err error) {
	sqlStr := `update videos
	set favorite_count=favorite_count+? 
	where id=?`
	_, err = db.Exec(sqlStr, num, videoId)
	if err != nil {
		zap.L().Error("mysql/AddFavoriteCount failed", zap.Error(err))
		return err
	}
	return nil
}

// QueryFavoriteVideos 查找用户喜欢视频的id
func QueryFavoriteVideos(userId int64) ([]int64, error) {
	sqlStr := `select video_id from user_favorite_video where user_id=? and is_favorite=true`
	videoIdList := make([]int64, 0)
	if err := db.Select(&videoIdList, sqlStr, userId); err != nil {
		zap.L().Error("mysql/QueryFavoriteVideo select failed", zap.Error(err))
		return videoIdList, err
	}
	if len(videoIdList) == 0 {
		return videoIdList, errors.New("没有喜欢的作品")
	}
	return videoIdList, nil
}

// QueryVideos 根据视频idList查询videoList  批量查找
func QueryVideos(videoIdList []int64) ([]models.VideosTable, error) {
	videoTableList := make([]models.VideosTable, len(videoIdList))
	//动态填充id
	strIDs := make([]string, 0, len(videoIdList))
	for _, id := range videoIdList {
		strIDs = append(strIDs, fmt.Sprintf("%d", id))
	}
	query, args, err := sqlx.In(`select id,author_id,play_url,cover_url,publish_time,favorite_count,comment_count,title
										from videos
										where id in(?)
										ORDER BY FIND_IN_SET(id, ?)`, videoIdList, strings.Join(strIDs, ","))
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	if err = db.Select(&videoTableList, query, args...); err != nil {
		zap.L().Error("mysql/favorite QueryVideos failed", zap.Error(err))
		return videoTableList, err
	}
	return videoTableList, nil
}

// AddUserFavoriteCount 用户点赞视频数量字段+1
func AddUserFavoriteCount(userId, num int64) error {
	sqlStr := `update user_show set favorite_count=favorite_count+? where user_id=?`
	_, err := db.Exec(sqlStr, num, userId)
	if err != nil {
		zap.L().Error("mysql/favorite AddUserFavoriteCount failed", zap.Error(err))
		return err
	}
	return nil
}

// AddAllFavoriteCount 视频用户获赞总数+1
func AddAllFavoriteCount(videoUserId, num int64) error {
	sqlStr := `update user_info set praise_num=praise_num+? where user_id=?`
	_, err := db.Exec(sqlStr, num, videoUserId)
	if err != nil {
		zap.L().Error("mysql/favorite AddAllFavoriteCount failed", zap.Error(err))
		return err
	}
	return nil
}
