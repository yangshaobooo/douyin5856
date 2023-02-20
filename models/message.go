package models

import "time"

type Messages struct {
	Id         int64     `json:"id" db:"id"`
	ToUserId   int64     `json:"to_user_id" db:"to_user_id"`
	FromUserId int64     `json:"from_user_id" db:"from_user_id"`
	Content    string    `json:"content" db:"content"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
}

// Message 要求消息的返回格式
type Message struct {
	Id         int64  `json:"id"`
	ToUserId   int64  `json:"to_user_id"`
	FromUserId int64  `json:"from_user_id"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
}

type FriendUser struct {
	User
	Avatar  string `json:"avatar"`
	Content string `json:"message,omitempty"`
	MsgType int64  `json:"msgType"`
}

// ResponseFriend 好友列表响应
type ResponseFriend struct {
	Response
	FriendList []FriendUser `json:"user_list"`
}

// RequestSend 发送消息的请求
type RequestSend struct {
	Token      string `form:"token"`
	ToUserId   int64  `form:"to_user_id"`
	ActionType int32  `form:"action_type"`
	Content    string `form:"content"`
}

// ResponseChatRecord 消息记录的响应
type ResponseChatRecord struct {
	Response
	MessageList []Message `json:"message_list"`
}
