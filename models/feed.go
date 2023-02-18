package models

import "time"

type VideosTable struct {
	VideoID       int64     `db:"id"`
	AuthorID      int64     `db:"author_id"`
	PlayUrl       string    `db:"play_url"`
	CoverUrl      string    `db:"cover_url"`
	PublishTime   time.Time `db:"publish_time"`
	FavoriteCount int64     `db:"favorite_count"`
	CommentCount  int64     `db:"comment_count"`
	Title         string    `db:"title"`
}
