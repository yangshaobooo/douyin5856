package logic

import (
	"douyin5856/dao/mysql"
	"douyin5856/models"
	"fmt"
	"go.uber.org/zap"
	"log"
	"sort"
	"strconv"
	"sync"
	"time"
)

// FriendList 处理朋友列表相关逻辑
func FriendList(userId int64) (*models.ResponseFriend, error) {
	// 1、根据用户id查询 user_follow表，找到所有的粉丝用户
	fansIdList, err := mysql.QueryFansID(userId)
	if err != nil {
		return nil, err
	}
	log.Printf("粉丝用户%v\n", fansIdList)
	// 2、查找所有的关注用户
	followIdList, err := mysql.QueryFollowID(userId)
	if err != nil {
		return nil, err
	}
	log.Printf("关注用户%v\n", followIdList)
	// 3、取交集
	friendIdList := InterFollowFans(fansIdList, followIdList)
	log.Printf("朋友%v\n", friendIdList)
	if len(friendIdList) == 0 {
		return nil, nil
	}
	// 4、创建朋友列表的切片
	friends := make([]models.FriendUser, len(friendIdList))
	wgFriendList := &sync.WaitGroup{}
	wgFriendList.Add(len(friendIdList))
	num := 0
	for _, friendId := range friendIdList {
		var friendInfo models.FriendUser
		go func(friendId int64) {
			StuffOneFriend(&friendInfo, friendId, userId)
			friends[num] = friendInfo
			num = num + 1
			wgFriendList.Done()
		}(friendId)
	}
	wgFriendList.Wait()
	// 同样好友列表也不用排序
	// 5、信息组合
	responseFriend := &models.ResponseFriend{
		models.Response{
			0,
			"好友列表成功",
		},
		friends,
	}
	return responseFriend, nil
}

// MessageAction 处理消息发送逻辑
func MessageAction(p *models.RequestSend, curUserId int64) (err error) {
	// 1、把消息存进数据库中
	messageTime := time.Now()
	if err = mysql.MessageAction(p.ToUserId, curUserId, p.Content, messageTime); err != nil {
		return err
	}
	// 2、消息存储成功
	return nil
}

// MessageChat 处理消息记录相关逻辑
func MessageChat(curUserId, toUserId int64) (*models.ResponseChatRecord, error) {
	// 1、根据发送者id和接收者id查询所有的消息
	messagesList, err := mysql.MessageChat(curUserId, toUserId)
	if err != nil {
		return nil, err
	}
	// 数据库中的格式转化成需要的格式
	toMessageList := make([]models.Message, len(messagesList))
	for i, p := range messagesList {
		toMessageList[i].Id = p.Id
		toMessageList[i].FromUserId = p.FromUserId
		toMessageList[i].ToUserId = p.ToUserId
		toMessageList[i].Content = p.Content
		toMessageList[i].CreateTime = strconv.FormatInt(p.CreateTime.UnixMilli(), 10)
		//toMessageList[i].CreateTime = p.CreateTime.Format("03:04 PM")
	}
	fmt.Println(toMessageList)
	// 2、组装返回的响应
	responChatRecord := &models.ResponseChatRecord{
		models.Response{
			0,
			"返回消息记录成功",
		},
		toMessageList,
	}
	return responChatRecord, nil
}

// StuffOneFriend 填充一条朋友信息
func StuffOneFriend(friendInfo *models.FriendUser, fansId, curUserId int64) {
	user, err := GetUserByIdWithCurId(fansId, curUserId)
	if err != nil {
		zap.L().Error("message StuffOneFriend failed", zap.Error(err))
	}
	friendInfo.User = user
	friendInfo.Avatar = "http://192.168.0.112:8080/static/touxiang.jpg"
	friendInfo.Content, friendInfo.MsgType, _ = mysql.QueryLatestChat(curUserId, fansId)
}

// InterFollowFans 去交集得到朋友列表
func InterFollowFans(fansIdList, followIdList []int64) []int64 {
	// 还是要排序的对fansIdList
	sort.Slice(fansIdList, func(i, j int) bool {
		return fansIdList[i] < fansIdList[j]
	})
	sort.Slice(followIdList, func(i, j int) bool {
		return followIdList[i] < followIdList[j]
	})
	friendIdList := make([]int64, 0)
	var (
		fansIndex   int
		followIndex int
	)
	// 时间复杂度O（n）。
	for fansIndex < len(fansIdList) && followIndex < len(followIdList) {
		if fansIdList[fansIndex] == followIdList[followIndex] {
			friendIdList = append(friendIdList, fansIdList[fansIndex])
			fansIndex++
			followIndex++
		} else if fansIdList[fansIndex] < followIdList[followIndex] {
			fansIndex++
		} else {
			followIndex++
		}
	}
	return friendIdList
}
