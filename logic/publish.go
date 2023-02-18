package logic

import (
	"douyin5856/dao/mysql"
	"douyin5856/middlewares/ffmpeg"
	"douyin5856/middlewares/snowflake1"
	"douyin5856/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"mime/multipart"
	"sort"
	"strconv"
	"sync"
	"time"
)

var (
	defaultVideoSuffix = ".mp4"
	defaultImageSuffix = ".jpg"
)

// Publish 处理视频上传相关逻辑
func Publish(c *gin.Context, userId int64, title string, data *multipart.FileHeader) (err error) {
	log.Println("Publish : running")
	// 1、给视频生成唯一ID
	videoID := snowflake1.GenID()
	wgPublish := &sync.WaitGroup{}
	wgPublish.Add(2)

	// 2、把视频保存在本地
	saveFile := "./public/" + strconv.FormatInt(videoID, 10) + ".mp4"
	fmt.Printf("保存的地址是:%v\n", saveFile)
	go func() {
		if err = c.SaveUploadedFile(data, saveFile); err != nil {
			zap.L().Error("保存视频到本地失败", zap.Error(err))
		}
		wgPublish.Done()
	}()
	// 3、给视频生成cover并且保存在本地
	coverFile := "./public/" + strconv.FormatInt(videoID, 10)
	go func() {
		_, err = ffmpeg.GetSnapshot(saveFile, coverFile, 1)
		if err != nil {
			zap.L().Error("生成图片失败", zap.Error(err))
		}
		wgPublish.Done()
	}()
	wgPublish.Wait()
	// 4、把视频相关数据存到mysql:视频id 用户id 播放地址 cover地址 发布时间 （点赞数 评论数默认0）title
	playUrl := "http://192.168.0.112:8080/static/" + strconv.FormatInt(videoID, 10) + ".mp4"
	coverUrl := "http://192.168.0.112:8080/static/" + strconv.FormatInt(videoID, 10) + ".jpg"
	fmt.Println(playUrl)
	fmt.Println(coverUrl)
	video := &models.VideosTable{
		VideoID:       videoID,
		AuthorID:      userId,
		PlayUrl:       playUrl,
		CoverUrl:      coverUrl,
		PublishTime:   time.Now(),
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
	}
	err = mysql.Publish(video)
	return
}

// PublishList 用户发布视频列表逻辑处理
func PublishList(userId, curId int64) (*models.ResponsePublishList, error) {
	// 1、根据用户id查找用户的视频
	videosTables, err := mysql.PublishList(userId)
	if err != nil {
		return nil, err
	}
	// 2、创建视频列表数组
	responseVideos := make([]models.Video, len(videosTables))
	// 3、并发查询视频数据
	wgPublishList := &sync.WaitGroup{}
	wgPublishList.Add(len(videosTables))
	num := 0
	for _, videosTable := range videosTables {
		var responseVideo models.Video
		go func(videoTable models.VideosTable) {
			// 填充一条视频数据
			StuffOneVideo(&responseVideo, &videoTable, curId)
			responseVideos[num] = responseVideo
			num = num + 1
			wgPublishList.Done()
		}(videosTable)
	}
	wgPublishList.Wait()
	// 同样按照id来排序，带有时间顺序的
	sort.Sort(FeedSlice(responseVideos))
	// 4、组装最后的返回数据
	responsePublishList := &models.ResponsePublishList{
		models.Response{
			0,
			"请求视频成功",
		},
		responseVideos,
	}
	return responsePublishList, nil
}
