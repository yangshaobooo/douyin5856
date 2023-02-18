package models

// Response 统一的响应码
type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}
type User struct {
	UserID        int64  `json:"id"`
	Username      string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

// ResponseSignUp 注册响应
type ResponseSignUp struct {
	Response
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

// ResponseLogin 登录响应
type ResponseLogin struct {
	Response
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

// ResponseUserInfo 用户信息响应
type ResponseUserInfo struct {
	Response
	User `json:"user"`
}

// ResponseFeed 视频流响应
type ResponseFeed struct {
	Response
	VideoList []Video `json:"video_list"`
	NextTime  int64   `json:"next_time"`
}

// Video 视频属性
type Video struct {
	ID            int64 `json:"id"`
	User          `json:"author"`
	PlayUrl       string `json:"play_url"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title"`
}

type ResponsePublishList struct {
	Response
	VideoList []Video `json:"video_list"`
}

type ResponseFavoriteList struct {
	Response
	VideoList []Video `json:"video_list"`
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

// 评论的响应
type ResponseComment struct {
	Response
	Comment `json:"comment"`
}

// 评论列表的响应
type ResponseCommentList struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

// 关注列表的响应
type ResponseUserList struct {
	Response
	UserList []User `json:"user_list"`
}
