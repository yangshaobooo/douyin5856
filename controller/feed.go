package controller

import (
	"douyin5856/logic"
	"douyin5856/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// Feed 视频流功能
func Feed(c *gin.Context) {
	//1、 获取请求参数
	p := new(models.RequestFeed)
	if err := c.ShouldBindQuery(p); err != nil { // 使用shouldbindQuery来接收
		// 请求参数有无
		zap.L().Error("Feed with invalid param", zap.Error(err))
		c.JSON(http.StatusOK, models.ResponseFeed{
			Response: models.Response{StatusCode: 1, StatusMsg: "获取视频流失败"},
		})
		return
	}
	zap.L().Info("传入的时间", zap.Int64("last_time", p.LatestTime))
	var lastTime time.Time
	if p.LatestTime != 0 {
		lastTime = time.Unix(p.LatestTime/1000, 0) //这里除1000我是真没想到，离大普
	} else {
		lastTime = time.Now()
	}
	zap.L().Info("获取到的时间戳", zap.Time("last_time", lastTime))
	userId := c.GetInt64("userID") //得不到默认是0
	fmt.Printf("视频流当前的用户名%v\n", userId)

	//2、处理视频feed流相关逻辑
	feed, err := logic.Feed(lastTime, userId) // 需要转化成int64

	if err != nil {
		c.JSON(http.StatusOK, models.ResponseFeed{
			Response: models.Response{StatusCode: 1, StatusMsg: "数据库中没有视频了"},
		})
		return
	}
	//fmt.Println(feed)
	c.JSON(http.StatusOK, feed)
}
