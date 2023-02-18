package mysql

import (
	"douyin5856/models"
	"go.uber.org/zap"
	"time"
)

// QueryLatestChat 查找用户的最新的那条聊天
func QueryLatestChat(userId, friendId int64) (chatContent string, messageType int64, err error) {
	// 查找两个用户的聊天记录
	list, err := MessageChat(userId, friendId)
	if err != nil {
		return "", 0, err
	}
	if len(list) == 0 {
		// 如果时没有聊天记录
		chatContent = ""
		messageType = 2
	} else {
		chatContent = list[len(list)-1].Content //我们要最后一条，因为之前是升序排序的
		if list[len(list)-1].FromUserId == userId {
			messageType = 0
		} else {
			messageType = 1
		}
	}
	return
}

// MessageAction 存储用户发送的信息
func MessageAction(toUserId, curUserId int64, content string, messageTime time.Time) (err error) {
	sqlStr := `insert into messages(to_user_id,from_user_id,content,create_time)values(?,?,?,?)`
	_, err = db.Exec(sqlStr, toUserId, curUserId, content, messageTime)
	if err != nil {
		zap.L().Error("mysql/message/ MessageAction failed", zap.Error(err))
		return err
	}
	// 存储成功
	return nil
}

// MessageChat 根据两个id查找消息记录
func MessageChat(curUserId, toUserID int64) ([]models.Messages, error) {
	var messageList []models.Messages
	sqlStr := `select id,to_user_id,from_user_id,content,create_time from messages
			where (to_user_id=? and from_user_id=?) or(to_user_id=? and from_user_id=?)
			order by create_time asc`
	if err := db.Select(&messageList, sqlStr, toUserID, curUserId, curUserId, toUserID); err != nil {
		return nil, err
	}
	// 没有错误
	return messageList, nil
}
