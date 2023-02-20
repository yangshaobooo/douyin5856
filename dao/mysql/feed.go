package mysql

import (
	"douyin5856/models"
	"errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

// Feed 取出latest时间之前的规定数量视频
func Feed(latestTime time.Time) ([]models.VideosTable, error) {
	// 读取配置文件中规定的视频数量
	count := viper.GetInt("video.count")
	videos := make([]models.VideosTable, count)
	sqlStr := `select id,author_id, play_url, cover_url, publish_time,favorite_count,comment_count,title 
				from videos 
				where publish_time<? 
				order by publish_time desc
				limit ?`
	// 查询多行数据使用select
	if err := db.Select(&videos, sqlStr, latestTime, count); err != nil {
		zap.L().Error("mysql/feed get failed", zap.Error(err))
		return videos, err
	}
	// 查询成功
	if len(videos) == 0 {
		zap.L().Error("没有新的视频了")
		return videos, errors.New("没有新的视频了")
	}
	return videos, nil
}

// QueryUserIdByVideoId 找到该视频的用户id
func QueryUserIdByVideoId(videoId int64) (int64, error) {
	sqlStr := `select author_id from videos where id = ?`
	var userId int64
	err := db.Get(&userId, sqlStr, videoId)
	if err != nil {
		zap.L().Error("mysql/feed QueryUserIdByVideoId failed", zap.Error(err))
		return userId, err
	}
	return userId, nil
}
