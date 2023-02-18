package mysql

import (
	"database/sql"
	"errors"
	"go.uber.org/zap"
)

// HaveRelation 判断关系表里面有没有
func HaveRelation(curUserId, userId int64) (bool, error) {
	var count int
	sqlStr := `select count(id) from user_follow where user_id=? and follower_id=?`
	if err := db.Get(&count, sqlStr, curUserId, userId); err != nil {
		zap.L().Error("mysql/relation/haveRelation failed", zap.Error(err))
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

// AlterRelation 更改用户的关注关系
func AlterRelation(curUserId, followerId int64) error {
	sqlStr := `update user_follow set is_follow = if(is_follow,0,1) where user_id=? and follower_id=?`
	_, err := db.Exec(sqlStr, curUserId, followerId)
	if err != nil {
		zap.L().Error("更改用户关系失败", zap.Error(err))
		return err
	}
	// 更改关系成功
	return nil
}

// InsertRelation 插入一段新的关注关系
func InsertRelation(curUserId, followerId int64, isFollow bool) error {
	sqlStr := `insert into user_follow(user_id,follower_id,is_follow)values(?,?,?)`
	_, err := db.Exec(sqlStr, curUserId, followerId, isFollow)
	if err != nil {
		zap.L().Error("mysql/relation/InsertRelation failed", zap.Error(err))
		return err
	}
	return nil
}

// FollowNumChange 关注数量改变
func FollowNumChange(userId, num int64) error {
	sqlStr := `update user_info set follow_num=follow_num+? where user_id =?`
	_, err := db.Exec(sqlStr, num, userId)
	if err != nil {
		zap.L().Error("mysql/relation FollowNumChange failed", zap.Error(err))
		return err
	}
	return nil
}

// FansNumChange 粉丝数量改变
func FansNumChange(userId, num int64) error {
	sqlStr := `update user_info set fans_num=fans_num+? where user_id =?`
	_, err := db.Exec(sqlStr, num, userId)
	if err != nil {
		zap.L().Error("mysql/relation FansNumChange failed", zap.Error(err))
		return err
	}
	return nil
}

// QueryFollowID 查询关注人的id
func QueryFollowID(userId int64) ([]int64, error) {
	sqlStr := `select follower_id from user_follow where user_id=? and is_follow=true`
	var followIdList []int64
	if err := db.Select(&followIdList, sqlStr, userId); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Debug("没有关注的的人")
			return followIdList, nil
		}
		zap.L().Error("查询关注的人失败", zap.Error(err))
		return nil, err
	}
	return followIdList, nil
}

// QueryFansID 查询粉丝列表
func QueryFansID(userId int64) ([]int64, error) {
	sqlStr := `select user_id from user_follow where follower_id=? and is_follow=true`
	var fansIdList []int64
	if err := db.Select(&fansIdList, sqlStr, userId); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Debug("没有粉丝的人")
			return fansIdList, nil
		}
		zap.L().Error("查询粉丝失败", zap.Error(err))
		return nil, err
	}
	if len(fansIdList) == 0 {
		return fansIdList, errors.New("没有粉丝")
	}
	return fansIdList, nil
}
