package models

// RequestSignUp 用户注册请求的结构体
type RequestSignUp struct {
	Username string `form:"username"`
	Password string `form:"password"` //这里前端传的是query 这里需要添加form标签
}

// RequestLogin 用户登录请求的结构体，虽然一样内容，但是为了以后的可扩展性，还是用不一样的
type RequestLogin struct {
	Username string `form:"username"`
	Password string `form:"password"` //这里前端传的是query 这里需要添加form标签
}

// RequestUserInfo 获取用户信息的请求
type RequestUserInfo struct {
	UserID int64  `form:"user_id"`
	Token  string `form:"token"`
}

// RequestFeed 视频流请求
type RequestFeed struct {
	LatestTime int64  `form:"latest_time"`
	Token      string `form:"token"`
}

// RequestPublish 发布视频请求
type RequestPublish struct {
	Token string `form:"token"`
	Data  []byte `form:"data"`
	Title string `form:"title"`
}

// RequestFavorite 视频点赞和取消点赞请求
type RequestFavorite struct {
	Token      string `form:"token"`
	VideoId    int64  `form:"video_id"`
	ActionType int32  `form:"action_type"`
}

// RequestCommentAction 视频评论操作
type RequestCommentAction struct {
	Token       string `form:"token"`
	VideoId     int64  `form:"video_id"`
	ActionType  int32  `form:"action_type"`
	CommentText string `form:"comment_text"`
	CommentId   int64  `form:"comment_id"`
}

// RequestRelation 关注操作
type RequestRelation struct {
	Token      string `form:"token"`
	ToUserId   int64  `form:"to_user_id"`
	ActionType int32  `form:"action_type"`
}
