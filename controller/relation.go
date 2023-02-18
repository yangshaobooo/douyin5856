package controller

import (
	"douyin5856/logic"
	"douyin5856/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// RelationAction 用户关注或者取消关注
func RelationAction(c *gin.Context) {
	// 1、获取参数
	p := new(models.RequestRelation)
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("获取参数失败", zap.Error(err))
		return
	}
	curUserId := c.GetInt64("userID")
	// 2、处理关注相关逻辑
	if err := logic.RelationAction(p, curUserId); err != nil {
		return
	}
	// 3、返回响应
	c.JSON(http.StatusOK, models.Response{
		0,
		"点赞成功",
	})
}

// FollowList 查看关注列表
func FollowList(c *gin.Context) {
	// 1、获取参数
	userIdS, _ := c.GetQuery("user_id")
	userId, _ := strconv.ParseInt(userIdS, 10, 64)

	// 2、处理查看关注列表相关逻辑
	responseUserList, err := logic.FollowList(userId)
	if err != nil {
		return
	}
	// 3、返回响应
	c.JSON(http.StatusOK, responseUserList)
}

// FansList 查询粉丝列表
func FansList(c *gin.Context) {
	// 1、获取参数
	userIdS, _ := c.GetQuery("user_id")
	userId, _ := strconv.ParseInt(userIdS, 10, 64)

	// 2、处理查看关注列表相关逻辑
	responseUserList, err := logic.FansList(userId)
	if err != nil {
		return
	}
	// 3、返回响应
	c.JSON(http.StatusOK, responseUserList)
}

// FriendList 好友列表
func FriendList(c *gin.Context) {
	// 1、处理参数
	userIdS, _ := c.GetQuery("user_id")
	userId, _ := strconv.ParseInt(userIdS, 10, 64)
	// 2、处理查看好友列表相关逻辑
	responseFriend, err := logic.FriendList(userId)
	if err != nil {
		return
	}
	// 3、返回响应
	c.JSON(http.StatusOK, responseFriend)
}
