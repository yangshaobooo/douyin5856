package logic

import (
	"douyin5856/dao/mysql"
	"douyin5856/middlewares/snowflake1"
	"douyin5856/models"
	"go.uber.org/zap"
	"log"
	"sort"
	"sync"
	"time"
)

// CommentAction 处理评论相关逻辑
func CommentAction(p *models.RequestCommentAction, userId int64) (*models.ResponseComment, error) {
	if p.ActionType == 2 {
		//删除评论
		// 从数据库中删除该条评论数据
		if err := mysql.CommentDelete(p.CommentId); err != nil {
			return nil, err
		}
		// 修改视频的评论数量-1
		if err := mysql.VideoCommentAdd(p.VideoId, -1); err != nil {
			return nil, err
		}
		return nil, nil
	}

	// 发布评论
	cTime := time.Now()
	// 1、雪花算法生成一个评论id
	commentId := snowflake1.GenID()
	if err := mysql.CommentStore(commentId, p.VideoId, userId, p.CommentText, cTime); err != nil {
		return nil, err
	}
	// 2、查询用户相关信息
	user, err := GetUserByIdWithCurId(userId, userId)
	if err != nil {
		return nil, err
	}
	// 2、视频的评论数量要+1
	if err = mysql.VideoCommentAdd(p.VideoId, 1); err != nil {
		return nil, err
	}
	sTime := cTime.Format("01-02")
	//4、组装需要返回的响应
	res := &models.ResponseComment{
		models.Response{
			0,
			"评论成功",
		},
		models.Comment{
			commentId,
			user,
			p.CommentText,
			sTime,
		},
	}
	return res, nil
}

// CommentList 处理评论列表函数
func CommentList(videoId, curUserId int64) (*models.ResponseCommentList, error) {
	log.Println("logic/comment/CommentList :running")
	//1、分局videoId 去comments里面查找所有的评论
	comments, err := mysql.QueryComments(videoId)
	if err != nil {
		return nil, err
	}
	// 2、定义好返回的评论列表
	resComment := make([]models.Comment, len(comments))
	// 3、并发提高查询速度  并发组装评论列表
	wg := &sync.WaitGroup{}
	wg.Add(len(comments))
	num := 0
	for _, commentTable := range comments {
		var commentData models.Comment
		go func(commentTable models.CommentsTable) {
			StuffOneComment(&commentData, &commentTable, curUserId)
			resComment[num] = commentData
			num = num + 1
			wg.Done()
		}(commentTable)
	}
	wg.Wait()
	// 因为协程打乱了之前的顺序，需要重新排序，根据主键排序
	sort.Sort(CommentSlice(resComment))

	// 4、加上response
	responseCommentList := &models.ResponseCommentList{
		models.Response{
			0,
			"请求评论列表成功",
		},
		resComment,
	}
	return responseCommentList, nil
}

// StuffOneComment 填充一条评论数据
func StuffOneComment(comment *models.Comment, commentTable *models.CommentsTable, curUserId int64) {
	// 评论表中自带的数据直接赋值
	comment.Id = commentTable.Id
	comment.Content = commentTable.CommentText
	comment.CreateDate = commentTable.CreateTime.Format("01-02")
	// 填充一条评论缺少的是User的相关数据，查询返回
	var err error
	comment.User, err = GetUserByIdWithCurId(commentTable.UserId, curUserId)
	if err != nil {
		zap.L().Error("logic/comment/StuffComment GetUserByIdWithCurId failed", zap.Error(err))
	}
}

// CommentSlice 为了评论的排序
type CommentSlice []models.Comment

func (a CommentSlice) Len() int {
	return len(a)
}
func (a CommentSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a CommentSlice) Less(i, j int) bool {
	return a[i].Id > a[j].Id
}
