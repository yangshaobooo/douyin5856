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
	Praise    int64 `db:"praise_num"`
}
