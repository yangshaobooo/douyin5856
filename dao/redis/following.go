package redis

import (
	"strconv"
)

//从redis中获取是不是目标的粉丝
func IsFans(userId,targetId int64)(exist bool,err error){
	if exist, err = rdbFollowing.SIsMember(strconv.FormatInt(userId, 10), strconv.FormatInt(targetId, 10)).Result();err!=nil{
		
		return 
	}
	return 
}
//成为粉丝吧
func BeFans(userId,targetId int64)(error){
	if _, err := rdbFollowing.SAdd(strconv.FormatInt(userId, 10), strconv.FormatInt(targetId, 10)).Result();err!=nil{
		
		return err
	}
	return nil
}
//删除粉丝关系
func NoFans(userId,target int64)(error){
	if _, err := rdbFollowing.SRem(strconv.FormatInt(userId, 10), strconv.FormatInt(target, 10)).Result();err!=nil{
		
		return err
	}
	return nil
}
//获取关注数量
func GetFollowingCnt(userId int64)(Cnt int64,err error){
	if Cnt, err = rdbFollowing.SCard(strconv.FormatInt(userId, 10)).Result();err!=nil{
		return
	}
	return
}
//根据用户id返回所有关注id
func GetFollowingList(userId int64)([]int64,error){
	result, err := rdbFollowing.SMembers(strconv.FormatInt(userId, 10)).Result()
	if err!=nil{
		return []int64{},err
	}
	IdList:=make([]int64,len(result))
	for i,v:=range result{
		IdList[i],_=strconv.ParseInt(v,10,64)
	}
	return IdList,nil

}