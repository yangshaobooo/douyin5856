package models

import "time"

type CommentsTable struct {
	Id          int64     `db:"id"`
	UserId      int64     `db:"user_id"`
	VideoId     int64     `db:"video_id"`
	CommentText string    `db:"comment_text"`
	CreateTime  time.Time `db:"create_date"`
}
