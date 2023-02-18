package jwt

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthMiddleware 基于JWT的认证中间件
func AuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		//fmt.Println("开始正常鉴权")
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		//authHeader := c.Request.Header.Get("Authorization")
		tokenStr := c.Query("token")
		if tokenStr == "" {
			tokenStr = c.PostForm("token")
		}
		if tokenStr == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "没有token",
			})
			c.Abort()
			return
		}
		tokenStruck, err := ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		//c.Set(controller.CtxUserIDKey, mc.UserId)
		//fmt.Printf("用户id是%v\n", tokenStruck.UserId)
		c.Set("userID", tokenStruck.UserId)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}

// AuthWithoutLogin 未登录情况下,若携带token,则解析出用户id并放入context;若未携带,则放入用户id默认值0
func AuthWithoutLogin() gin.HandlerFunc {
	return func(context *gin.Context) {
		//fmt.Println("开始视频流的用户鉴权")
		auth := context.Query("token")
		var userId int64
		if len(auth) == 0 {
			userId = 0
			//fmt.Println("未登录状态")
		} else {
			//fmt.Println("登录状态")
			token, err := ParseToken(auth)
			if err != nil {
				context.Abort()
			} else {
				userId = token.UserId
			}
		}
		context.Set("userID", userId)
		context.Next()
	}
}
