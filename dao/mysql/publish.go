package mysql

import (
	"douyin5856/models"
	"go.uber.org/zap"
)

// Publish 视频数据添加到数据库中
func Publish(video *models.VideosTable) error {
	sqlStr := `insert into videos
	(id,author_id,play_url,cover_url,publish_time,favorite_count,comment_count,title)
	values(?,?,?,?,?,?,?,?)`
	_, err := db.Exec(sqlStr, video.VideoID, video.AuthorID,
		video.PlayUrl, video.CoverUrl, video.PublishTime,
		video.FavoriteCount, video.CommentCount, video.Title)
	if err != nil {
		zap.L().Error("视频存入数据库失败", zap.Error(err))
		return err
	}
	return nil
}

// PublishCountAdd 用户发布视频数量字段+1
func PublishCountAdd(userId int64) error {
	sqlStr := `update user_show set work_count=work_count+1 where user_id = ?`
	_, err := db.Exec(sqlStr, userId)
	if err != nil {
		zap.L().Error("mysql/publish PublishCountAdd db.exec failed", zap.Error(err))
		return err
	}
	return nil
}

// PublishList 返回发布视频列表
func PublishList(userId int64) ([]models.VideosTable, error) {
	var videosTableList []models.VideosTable
	sqlStr := `select
		id,author_id,play_url,cover_url,publish_time,favorite_count,comment_count,title
		from videos 
		where author_id=?`
	if err := db.Select(&videosTableList, sqlStr, userId); err != nil {
		zap.L().Error("mysql/PublishList failed", zap.Error(err))
		return videosTableList, err
	}
	return videosTableList, nil
}
