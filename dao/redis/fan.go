package redis

import "strconv"


//粉丝量+1
func AddFollower(userId,targetId int64)(error){
	if _, err := rdbFans.SAdd(strconv.FormatInt(userId, 10), strconv.FormatInt(targetId, 10)).Result();err!=nil{
		return err
	}

	return nil 
}
//粉丝量-1
func CancelFollow(userId,targetId int64)(error){
	if _, err := rdbFans.SRem(strconv.FormatInt(userId, 10), strconv.FormatInt(targetId, 10)).Result();err!=nil{
		return err
	}
	return nil
}
//获取粉丝数量
func GetFansCnt(userId int64 )(int64,error){
	Cnt, err := rdbFans.SCard(strconv.FormatInt(userId, 10)).Result()
	if err!=nil{
		return Cnt,err
	}
	return Cnt,nil
}
//
func GetFansList(userId int64)([]int64,error){
	result, err := rdbFans.SMembers(strconv.FormatInt(userId, 10)).Result()
	if err!=nil{
		return []int64{},err
	}
	FansList:=make([]int64,len(result))
	for i,v:=range result{
		FansList[i],_=strconv.ParseInt(v,10,64)
	}
	return FansList,nil
}