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

// Publish 发布视频
func Publish(c *gin.Context) {
	//1、首先获得请求参数，进行校验,解析出用户名
	userId := c.GetInt64("userID")  //得到的就是int64 没得到默认是0
	data, err := c.FormFile("data") //获取文件
	if err != nil {
		zap.L().Error("获取不到视频", zap.Error(err))
		return
	}
	title := c.PostForm("title") // 获取标题

	// 2、处理发布视频相关逻辑
	err = logic.Publish(c, userId, title, data)
	if err != nil {
		return
	}
	//// 3、返回响应
	c.JSON(http.StatusOK, models.Response{
		StatusCode: 0,
		StatusMsg:  "upload success",
	})
}

// PublishList 用户发布视频列表
func PublishList(c *gin.Context) {
	// 1、获取参数
	userIds, _ := c.GetQuery("user_id")
	userId, _ := strconv.ParseInt(userIds, 10, 64) //string转化成int64
	curId := c.GetInt64("userID")                  // 获取当前用户id

	// 2、处理视频列表相关逻辑
	responsePublishList, err := logic.PublishList(userId, curId)
	if err != nil {
		log.Printf("调用logic.PublishList(%v)出现错误：%v\n", userId, err)
		c.JSON(http.StatusOK, models.ResponsePublishList{
			Response: models.Response{StatusCode: 1, StatusMsg: "获取视频列表失败"},
		})
		return
	}
	// 3、返回响应
	c.JSON(http.StatusOK, responsePublishList)
}
