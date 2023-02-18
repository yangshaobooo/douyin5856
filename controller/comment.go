package controller

import (
	"douyin5856/logic"
	"douyin5856/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
)

// CommentAction 用户评论功能
func CommentAction(c *gin.Context) {
	log.Println("Controller CommentAction start: running")
	// 1、参数校验功能
	p := new(models.RequestCommentAction)
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("请求参数错误", zap.Error(err))
		return
	}
	// 获取当前的用户id
	userId := c.GetInt64("userID")
	log.Printf("comment userId is %d", userId)
	// 2、处理评论的相关逻辑
	responseComment, err := logic.CommentAction(p, userId)
	if err != nil {
		return
	}
	log.Println("comment success")
	// 3、返回响应
	c.JSON(http.StatusOK, responseComment)
}

// CommentList 获取评论列表
func CommentList(c *gin.Context) {
	// 1、获取参数
	videoIdS, _ := c.GetQuery("video_id")
	videoId, _ := strconv.ParseInt(videoIdS, 10, 64)
	userId := c.GetInt64("userID")

	//2、处理评论列表相关逻辑
	responseCommentList, err := logic.CommentList(videoId, userId)
	if err != nil {
		c.JSON(http.StatusOK, models.ResponseCommentList{
			Response: models.Response{
				0,
				"当前视频没有评论",
			},
		})
		return
	}

	// 3、返回响应
	c.JSON(http.StatusOK, responseCommentList)
}
