package controller

import (
	"douyin5856/logic"
	"douyin5856/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// Register 用户注册函数
func Register(c *gin.Context) {
	// 1、------获取前端输入的username 和password------
	p := new(models.RequestSignUp)               //使用预先定义好的结构体来接受数据
	if err := c.ShouldBindQuery(p); err != nil { // 使用shouldBindQuery来接收 千万别用shouldbindjson
		// 请求参数有无
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		return
	}

	// 2、-----处理注册逻辑------
	err, res := logic.SignUp(p)
	if err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
	} else {
		zap.L().Info("用户注册成功")
	}

	// 3、-----返回响应------
	c.JSON(http.StatusOK, res)

}

// Login 用户登录函数
func Login(c *gin.Context) {
	// 1、------获取前端输入的username 和password------
	p := new(models.RequestLogin)                //使用预先定义好的结构体来接受数据
	if err := c.ShouldBindQuery(p); err != nil { // 使用shouldBindQuery来接收
		// 请求参数有误
		zap.L().Error("Login with invalid param", zap.Error(err))
		return
	}
	fmt.Println("登录的用户", p)

	// 2、------处理登录的逻辑------
	err, res := logic.Login(p)
	if err != nil {
		zap.L().Error("Username or password error", zap.Error(err))
		return
	} else {
		zap.L().Info("用户登录成功", zap.String("username", p.Username))
	}

	// 3、返回响应
	c.JSON(http.StatusOK, res)
}

// UserInfo 用户信息函数
func UserInfo(c *gin.Context) {
	// 1、获取query 包括user_id 和token
	p := new(models.RequestUserInfo)
	if err := c.ShouldBindQuery(p); err != nil { // 使用shouldBindQuery来接收
		// 请求参数有无
		zap.L().Error("UserInfo with invalid param", zap.Error(err))
		return
	}
	curUserId := c.GetInt64("userID")
	// 2、处理获取用户信息逻辑
	err, res := logic.UserInfo(p, curUserId)
	if err != nil {
		zap.L().Error("get userinfo failed", zap.Error(err))
		return
	} else {
		zap.L().Info("获取用户数据成功", zap.String("username:", res.Username))
	}
	// 3、返回响应
	c.JSON(http.StatusOK, res)
}
