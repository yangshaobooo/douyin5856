package routes

import (
	"douyin5856/controller"
	"douyin5856/logger"
	"douyin5856/middlewares/jwt"
	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", jwt.AuthWithoutLogin(), controller.Feed)
	apiRouter.GET("/user/", jwt.AuthMiddleware(), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", jwt.AuthMiddleware(), controller.Publish)
	apiRouter.GET("/publish/list/", jwt.AuthMiddleware(), controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", jwt.AuthMiddleware(), controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", jwt.AuthMiddleware(), controller.FavoriteList)
	apiRouter.POST("/comment/action/", jwt.AuthMiddleware(), controller.CommentAction)
	apiRouter.GET("/comment/list/", jwt.AuthMiddleware(), controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", jwt.AuthMiddleware(), controller.RelationAction) // 关注操作
	apiRouter.GET("/relation/follow/list/", jwt.AuthMiddleware(), controller.FollowList) // 关注列表
	apiRouter.GET("/relation/follower/list/", jwt.AuthMiddleware(), controller.FansList)
	apiRouter.GET("/relation/friend/list/", jwt.AuthMiddleware(), controller.FriendList)
	apiRouter.GET("/message/chat/", jwt.AuthMiddleware(), controller.MessageChat)      //消息记录
	apiRouter.POST("/message/action/", jwt.AuthMiddleware(), controller.MessageAction) // 发送消息

	return r
}
