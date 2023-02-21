package redis

import (
	"log"
	"strconv"
)
//给视频id插入评论id
func AddCommentById(videoId, CommentId int64) error {
	if _, err := rdbComment.SAdd(strconv.FormatInt(videoId, 10), strconv.FormatInt(CommentId, 10)).Result();err!=nil{
		log.Println(err.Error())
		return err
	}
	log.Println("redis addcommentbyId SUCCEED")
	return nil
}
//根据视频id获取命令数量
func GetCommentCntById(videoId int64)(int64,error){
	var Cnt int64
	var err error
	if Cnt, err = rdbComment.SCard(strconv.FormatInt(videoId, 10)).Result();err!=nil{
		log.Println(err.Error())
		return 0,err
	}

	return Cnt,nil
}
//删除评论
func DelCommentById(videoId,commentId int64)(error){
	if _, err := rdbComment.SRem(strconv.FormatInt(videoId, 10), strconv.FormatInt(commentId, 10)).Result();err!=nil{
		log.Println(err.Error())
		return err
	}
	return nil

}