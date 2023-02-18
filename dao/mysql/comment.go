package mysql

import (
	"douyin5856/models"
	"errors"
	"go.uber.org/zap"
	"time"
)

// CommentStore 评论存储功能
func CommentStore(commentId, videoId, userId int64, commentText string, cTime time.Time) error {
	sqlStr := `insert into comments(id,user_id,video_id,comment_text,create_date)
			values(?,?,?,?,?)`
	_, err := db.Exec(sqlStr, commentId, userId, videoId, commentText, cTime)
	if err != nil {
		zap.L().Error("评论存储失败", zap.Error(err))
		return err
	}
	// 评论存储成功
	return nil
}

// CommentDelete 评论删除
func CommentDelete(commentId int64) error {
	sqlStr := `delete from comments where id=?`
	_, err := db.Exec(sqlStr, commentId)
	if err != nil {
		zap.L().Error("删除评论失败", zap.Error(err))
		return err
	}
	// 删除评论成功
	return nil
}

// VideoCommentAdd 视频的评论数量需要加num
func VideoCommentAdd(videoId int64, num int) error {
	sqlStr := `update videos set comment_count=comment_count+? where id = ?`
	_, err := db.Exec(sqlStr, num, videoId)
	if err != nil {
		zap.L().Error("mysql/comment/VideoCommentAdd failed", zap.Error(err))
		return err
	}
	return nil
}

// QueryComments 查找视频的所有评论
func QueryComments(videoId int64) ([]models.CommentsTable, error) {
	var commentList []models.CommentsTable
	sqlStr := `select id,user_id,video_id,comment_text,create_date from comments where video_id=? order by create_date desc `
	if err := db.Select(&commentList, sqlStr, videoId); err != nil {
		zap.L().Error("mysql /comment/QueryComments failed", zap.Error(err))
		return nil, err
	}
	if len(commentList) == 0 {
		return commentList, errors.New("当前视频没有评论")
	}
	return commentList, nil
}
