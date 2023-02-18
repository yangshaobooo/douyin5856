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

// FavoriteAction 用户点赞功能
func FavoriteAction(c *gin.Context) {
	// 1、获取请求参数
	p := new(models.RequestFavorite)
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("controller/FavoriteAction get request failed", zap.Error(err))
		return
	}
	userID := c.GetInt64("userID")
	// 2、处理点赞相关逻辑
	if err := logic.FavoriteAction(p.VideoId, userID, p.ActionType); err != nil {
		return
	}
	// 3、返回响应
	c.JSON(http.StatusOK, models.Response{
		0,
		"点赞操作成功",
	})
}

// FavoriteList 用户喜欢列表功能
func FavoriteList(c *gin.Context) {
	// 1、获取请求参数
	userIdS, _ := c.GetQuery("user_id")
	userId, _ := strconv.ParseInt(userIdS, 10, 64)

	curUserId := c.GetInt64("userID")

	// 2、处理喜欢列表相关逻辑
	responseFavoriteList, err := logic.FavoriteList(userId, curUserId)
	if err != nil {
		log.Printf("方法logic.FavoriteList 失败：%v", err)
		c.JSON(http.StatusOK, models.ResponseFavoriteList{
			Response: models.Response{
				StatusCode: 1,
				StatusMsg:  "get favouriteList fail ",
			},
		})
		return
	}
	// 3、返回响应
	c.JSON(http.StatusOK, responseFavoriteList)
}
