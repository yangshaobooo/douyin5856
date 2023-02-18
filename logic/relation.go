package logic

import (
	"douyin5856/dao/mysql"
	"douyin5856/models"
	"go.uber.org/zap"
	"log"
	"sync"
)

// RelationAction 处理用户的关注和取消关注
func RelationAction(p *models.RequestRelation, curUserId int64) (err error) {
	// 1、首先查询关系表里面有没有
	have, err := mysql.HaveRelation(curUserId, p.ToUserId)
	if err != nil {
		return
	}
	if have {
		//曾经有过关系，直接修改最后的is_follow就可以了
		if err = mysql.AlterRelation(curUserId, p.ToUserId); err != nil {
			return err
		}
	} else {
		// 需要向表中插入一段新的关系，必然是关注的操作
		if err = mysql.InsertRelation(curUserId, p.ToUserId, true); err != nil {
			return err
		}
	}
	if p.ActionType == 1 {
		// 当前用户的关注数量+1
		if err = mysql.FollowNumChange(curUserId, 1); err != nil {
			return err
		}
		// 目标用户的粉丝数量+1
		if err = mysql.FansNumChange(p.ToUserId, 1); err != nil {
			return err
		}
	} else {
		// 当前用户的关注数量-1
		if err = mysql.FollowNumChange(curUserId, -1); err != nil {
			return err
		}
		// 目标用户的粉丝数量-1
		if err = mysql.FansNumChange(p.ToUserId, -1); err != nil {
			return err
		}
	}
	return nil
}

// FollowList 查看关注列表逻辑
func FollowList(userId int64) (*models.ResponseUserList, error) {
	log.Println("FollowList : running")
	// 1、根据用户id查询 user_follow表，找到所有的关注用户
	followIdList, err := mysql.QueryFollowID(userId)
	if err != nil {
		return nil, err
	}
	// 2、创建用户关注的人的信息切片
	followsInfo := make([]models.User, len(followIdList))
	// 3、并发填充用户数据
	wgFollowList := &sync.WaitGroup{}
	wgFollowList.Add(len(followIdList))
	num := 0
	for _, followId := range followIdList {
		go func(followId int64) {
			user, err := GetUserByIdWithCurId(followId, userId)
			if err != nil {
				zap.L().Error("填充用户信息错误")
			}
			followsInfo[num] = user
			num = num + 1
			wgFollowList.Done()
		}(followId)
	}
	wgFollowList.Wait()
	// 没有必要排序了，随意就行
	// 4、信息组合
	responseUserList := &models.ResponseUserList{
		models.Response{
			0,
			"查询关注列表成功",
		},
		followsInfo,
	}
	return responseUserList, nil
}

// FansList 查询粉丝列表
func FansList(userId int64) (*models.ResponseUserList, error) {
	log.Println("FansList : running")
	// 1、根据用户id查询 user_follow表，找到所有的粉丝用户
	fansIdList, err := mysql.QueryFansID(userId)
	if err != nil {
		return nil, err
	}
	// 2、根据粉丝id填充粉丝的用户信息
	fansList := make([]models.User, len(fansIdList))
	wgFansList := &sync.WaitGroup{}
	wgFansList.Add(len(fansIdList))
	num := 0
	for _, fansId := range fansIdList {
		go func(fansId int64) {
			fansInfo, err := GetUserByIdWithCurId(fansId, userId)
			if err != nil {
				zap.L().Error("填充用户信息错误")
			}
			fansList[num] = fansInfo
			num = num + 1
			wgFansList.Done()
		}(fansId)
	}
	wgFansList.Wait()
	// 不需要排序了，随便顺序
	// 3、组合最后的数据
	responseUserList := &models.ResponseUserList{
		models.Response{
			0,
			"查询关注列表成功",
		},
		fansList,
	}
	return responseUserList, nil
}
