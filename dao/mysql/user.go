package mysql

import (
	"crypto/md5"
	"database/sql"
	"douyin5856/models"
	"encoding/hex"
	"errors"
	"go.uber.org/zap"
)

// 加密密码的密钥
const secret = "yang.com"

// CheckUserExist 判断用户在数据库中存不存在
func CheckUserExist(username string) (err error) {
	sqlStr := "select count(username) from user_basic where username=?"
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户已经存在")
	}
	return
}

// SignUp 想数据库中插入一条新的用户数据
func SignUp(user *models.UserBasic) (err error) {
	// 对密码进行加密  后面正式运行的时候打开就行  如需开启密码加密，取消注释即可
	//user.Password = encryptPassword(user.Password)

	// 执行sql语句
	sqlStr := `insert into user_basic(user_id,username,password)values(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	if err != nil {
		return
	}
	//创建用户的时候，把userInfo表一起创建
	sqlStr = `insert into user_info(user_id,follow_num,fans_num,praise_num)values (?,?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, 0, 0, 0)

	// 创建用户的收，把user_show一起创建了
	sqlStr = `insert into user_show(user_id,work_count,favorite_count)values (?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, 0, 0)
	return
}

func Login(user *models.UserBasic) (err error, userID int64) {
	oPassword := user.Password
	strSql := `select user_id,username,password from user_basic where username=?`
	// 查询该条数据
	if err := db.Get(user, strSql, user.Username); err != nil {
		// 没有查到用户
		if err == sql.ErrNoRows {
			return errors.New("该用户不存在"), 0
		}
		//查询数据库失败
		return err, 0
	}
	// ----------如果密码进行了加密，需要解密------------
	//encryptPassword(oPassword)

	//判断密码是否正确
	if oPassword != user.Password {
		return errors.New("密码错误"), 0
	}
	// 成功查询到用户,返回nil
	//需要userID,一起返回
	return nil, user.UserID
}

// encryptPassword 使用md5加密以下密码
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

// QueryName 查询表user根据userid查询username
func QueryName(userId int64) (string, error) {
	var userName string
	sqlStr := `select username from user_basic where user_id=?`
	if err := db.Get(&userName, sqlStr, userId); err != nil {
		if err == sql.ErrNoRows {
			return userName, errors.New("QueryName用户不存在")
		}
		zap.L().Error("mysql/QueryName failed", zap.Error(err))
		return userName, err
	}
	return userName, nil
}

// QueryFollowCount 查询表user_info 根据userid查询follow数量
func QueryFollowCount(userId int64) (int64, error) {
	var p int64
	sqlStr := `select follow_num from user_info where user_id=?`
	if err := db.Get(&p, sqlStr, userId); err != nil {
		if err == sql.ErrNoRows {
			return p, nil
		}
		return p, errors.New("获取用户失败粉丝和关注失败")
	}
	return p, nil
}

// QueryFansCount 查询表user_info 根据userid查询fans数量
func QueryFansCount(userId int64) (int64, error) {
	var p int64
	sqlStr := `select fans_num from user_info where user_id=?`
	if err := db.Get(&p, sqlStr, userId); err != nil {
		if err == sql.ErrNoRows {
			return p, nil
		}
		return p, errors.New("获取用户失败粉丝和关注失败")
	}
	return p, nil
}

// QueryUserInfo 查找用户的关注、粉丝、获赞 三个数据
func QueryUserInfo(userId int64) (*models.UserInfo, error) {
	p := new(models.UserInfo)
	sqlStr := `select user_id,follow_num,fans_num,praise_num from user_info where user_id=?`
	if err := db.Get(p, sqlStr, userId); err != nil {
		if err == sql.ErrNoRows {
			return p, nil
		}
		zap.L().Error("mysql/user QueryUserInfo db.get failed", zap.Error(err))
		return p, err
	}
	return p, nil
}

// UserRelation 查询表user_follow 是否有关注关系
func UserRelation(cur, des int64) (bool, error) {
	sqlStr := `select is_follow from user_follow use index(idx_user_follower)where user_id=? and follower_id=?`
	var isFollow bool
	if err := db.Get(&isFollow, sqlStr, cur, des); err != nil {
		if err == sql.ErrNoRows {
			return false, nil //表示一行没查到，表示两者之间没有关注关系
		}
		zap.L().Error("mysql UserRelation failed", zap.Error(err))
		return false, errors.New("mysql 查询用户关系失败")
	}
	// 成功查到
	return isFollow, nil
}

// QueryPubFavCount 查找发布和喜欢视频数量
func QueryPubFavCount(userId int64) (*models.UserShow, error) {
	sqlStr := `select user_id,work_count,favorite_count from user_show where user_id =?`
	p := new(models.UserShow)
	if err := db.Get(p, sqlStr, userId); err != nil {
		zap.L().Error("mysql/user QueryPubFavCount db.get failed", zap.Error(err))
		return p, err
	}
	return p, nil
}
