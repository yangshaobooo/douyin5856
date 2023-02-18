package controller

import (
	"douyin5856/logic"
	"douyin5856/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// MessageAction 发送消息功能
func MessageAction(c *gin.Context) {
	// 1、接受参数
	p := new(models.RequestSend)
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("MessageAction 接受参数出错", zap.Error(err))
		return
	}
	curUserId := c.GetInt64("userID")

	// 2、发送消息的逻辑处理
	if err := logic.MessageAction(p, curUserId); err != nil {
		return
	}
	// 3、 返回响应
	c.JSON(http.StatusOK, models.Response{
		0,
		"消息发送成功",
	})
}

// MessageChat 消息记录
func MessageChat(c *gin.Context) {
	// 1、 接受请求参数
	curUserId := c.GetInt64("userID")
	toUserIdString, _ := c.GetQuery("to_user_id")
	toUserId, _ := strconv.ParseInt(toUserIdString, 10, 64)

	// 2、消息记录的逻辑处理
	responseChatRecord, err := logic.MessageChat(curUserId, toUserId)
	if err != nil {
		return
	}
	// 3、返回响应
	c.JSON(http.StatusOK, responseChatRecord)
}
