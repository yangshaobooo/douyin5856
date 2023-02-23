package logic

import (
	"douyin5856/dao/mysql"
	"douyin5856/middlewares/jwt"
	"douyin5856/middlewares/snowflake1"
	"douyin5856/models"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"strconv"
	"sync"
)

// SignUp 注册逻辑处理
func SignUp(p *models.RequestSignUp) (error, *models.ResponseSignUp) {
	// 创建返回的response
	res := &models.ResponseSignUp{
		models.Response{
			-1,
			"",
		},
		0,
		"",
	}
	// 1、------判断用户在数据库中存不存在-------
	if err := mysql.CheckUserExist(p.Username); err != nil {
		zap.L().Error("用户已经存在", zap.Error(err))
		res.StatusMsg = "用户名已存在"
		return err, res
	}

	// 2、 使用雪花算法生成UID
	userID := snowflake1.GenID()

	// 3、把用户存进数据库
	user := models.UserBasic{
		userID,
		p.Username,
		p.Password,
	}
	if err := mysql.SignUp(&user); err != nil {
		zap.L().Error("mysql.InsertUser failed", zap.Error(err))
		return err, nil
	}

	// 4、生成token
	token, err := jwt.GenToken(userID, user.Username)
	if err != nil {
		zap.L().Error("jwt.GenToken failed", zap.Error(err))
		return err, nil
	}

	// 5、拼接注册最后的返回结构体
	res.StatusCode = 0
	res.StatusMsg = "注册成功"
	res.UserID = userID
	res.Token = token
	return err, res
}

// Login 登录逻辑处理
func Login(p *models.RequestLogin) (error, *models.ResponseLogin) {
	user := &models.UserBasic{
		Username: p.Username,
		Password: p.Password,
	}
	// 创建返回的response
	res := &models.ResponseLogin{
		models.Response{
			-1,
			"账号密码错误",
		},
		0,
		"",
	}

	// 查询数据库中用户名和密码
	err, userID := mysql.Login(user)
	if err != nil {
		return err, res
	}

	// 生成token
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return err, res
	}

	//拼接返回的response
	res.StatusCode = 0
	res.StatusMsg = "登录成功"
	res.UserID = userID
	res.Token = token
	return nil, res
}

// UserInfo 查看用户信息逻辑处理
func UserInfo(p *models.RequestUserInfo, curUserId int64) (error, *models.ResponseUserInfo) {
	log.Println("userInfo :running")
	if p.UserID == 0 {
		p.UserID = curUserId
	}
	user, err := GetUserByIdWithCurId(p.UserID, curUserId)
	if err != nil {
		return err, nil
	}
	res := &models.ResponseUserInfo{
		models.Response{
			0,
			"拉取用户数据成功",
		},
		user,
	}

	fmt.Println(res)
	return err, res
}

// GetUserByIdWithCurId 得到一个响应要求格式的User
func GetUserByIdWithCurId(userId, curUserId int64) (models.User, error) {
	// 需要返回的User
	user := models.User{
		UserID:        0,
		Username:      "",
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
		// 下面省略写，因为是可以省略的字段
	}
	user.UserID = userId
	user.Avatar = viper.GetString("app.address") + strconv.Itoa(viper.GetInt("app.port")) + "/static/" + "touxiang2.jpeg"          // 头像暂时都用一样的
	user.BackGroundImage = viper.GetString("app.address") + strconv.Itoa(viper.GetInt("app.port")) + "/static/" + "background.jpg" //背景图都用一样的
	user.Signature = "十年饮冰，难凉热血。"
	// 使用协程加快速度
	var wg sync.WaitGroup
	wg.Add(4)
	// 查询用户名字
	go func() {
		name, err := mysql.QueryName(userId)
		if err != nil {
			zap.L().Error("logic/user GetUserByIdWithCurId QueryName failed", zap.Error(err))
		}
		user.Username = name
		wg.Done()
	}()

	//// 查询用户的关注数量
	//go func() {
	//	// 数据库中查询
	//	followNum, err := mysql.QueryFollowCount(userId)
	//	if err != nil {
	//		zap.L().Error("logic/user GetUserByIdWithCurId QueryFollowCount failed", zap.Error(err))
	//	}
	//	user.FollowCount = followNum
	//	wg.Done()
	//}()
	//
	//// 查询用户的粉丝数量
	//go func() {
	//	// 数据库中查询
	//	fansNum, err := mysql.QueryFansCount(userId)
	//	if err != nil {
	//		zap.L().Error("logic/user GetUserByIdWithCurId QueryFansCount failed", zap.Error(err))
	//	}
	//	user.FollowerCount = fansNum
	//	wg.Done()
	//}()

	// 查询用户的关注、粉丝、获赞数量
	go func() {
		// 数据库中查询
		userInfo, err := mysql.QueryUserInfo(userId)
		if err != nil {
			zap.L().Error("logic/user GetUserByIdWithCurId mysql.QueryUserInfo failed", zap.Error(err))
		}
		user.FollowCount = userInfo.FollowNum
		user.FollowerCount = userInfo.FansNum
		user.TotalFavorited = userInfo.Praise
		wg.Done()
	}()

	// 查询当前用户是否关注该用户
	go func() {
		if curUserId == userId {
			user.IsFollow = true
		} else {
			relation, err := mysql.UserRelation(curUserId, userId)
			if err != nil {
				zap.L().Error("logic/user GetUserByIdWithCurId UserRelation failed", zap.Error(err))
			}
			user.IsFollow = relation
		}
		wg.Done()
	}()

	// 查询发布作品数量和点赞视频的数量
	go func() {
		userShow, err := mysql.QueryPubFavCount(userId)
		if err != nil {
			zap.L().Error("logic/user GetUserByIdWithCurId mysql.QueryPubFavCount failed", zap.Error(err))
		}
		user.WorkCount = userShow.WorkCount
		user.FavoriteCount = userShow.FavoriteCount
		wg.Done()
	}()
	wg.Wait()
	return user, nil
}
