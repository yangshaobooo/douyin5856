package models

type UserBasic struct {
	UserID   int64  `db:"user_id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

type UserInfo struct {
	UserID    int64 `db:"user_id"`
	FollowNum int64 `db:"follow_num"`
	FansNum   int64 `db:"fans_num"`
	Praise    int64 `db:"praise_num"` // 用户点赞总数量
}

type UserShow struct {
	UserID        int64 `db:"user_id"`
	WorkCount     int64 `db:"work_count"`     // 用户发布作品数量
	FavoriteCount int64 `db:"favorite_count"` // 用户喜欢视频的数量
}
