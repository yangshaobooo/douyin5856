package logic

import (
	"douyin5856/dao/mysql"
	"douyin5856/dao/redis"
	"douyin5856/models"
	"go.uber.org/zap"
	"log"
	"sort"
	"sync"
	"time"
)

// Feed 处理视频流逻辑
func Feed(latestTime time.Time, curUserId int64) (*models.ResponseFeed, error) {
	log.Println("feed : running")
	// 1、根据时间取出规定数量的视频信息
	res := new(models.ResponseFeed)
	videosList, err := mysql.Feed(latestTime)
	//fmt.Println(videosList)
	if err != nil {
		return res, err
	}
	// 2、组装视频列表 使用协程加快速度
	responseVideos := make([]models.Video, len(videosList))
	wgFeed := &sync.WaitGroup{}
	wgFeed.Add(len(videosList))
	num := 0
	for _, videoTable := range videosList {
		var responseVideo models.Video
		go func(videoTable models.VideosTable) {
			// 填充一条视频数据
			StuffOneVideo(&responseVideo, &videoTable, curUserId)
			responseVideos[num] = responseVideo
			num = num + 1
			wgFeed.Done()
		}(videoTable)
	}
	wgFeed.Wait()
	// 协程打乱了视频的顺序，因为不返回时间，我们使用videoId来排序
	// 因为使用了snowflake 雪花算法，videoId是带有时间属性的。
	sort.Sort(FeedSlice(responseVideos))

	// 4、组装返回的响应
	responseFeed := &models.ResponseFeed{
		models.Response{
			0,
			"请求视频成功",
		},
		responseVideos,
		videosList[len(videosList)-1].PublishTime.Unix() * 1000, // 这里需要成1000
	}
	return responseFeed, err
}

// StuffOneVideo 组装一条视频数据
func StuffOneVideo(responseVideo *models.Video, videoTable *models.VideosTable, curUserId int64) {
	wgStuffVideo := &sync.WaitGroup{}
	wgStuffVideo.Add(2)

	// 组装一条视频数据
	responseVideo.ID = videoTable.VideoID
	responseVideo.PlayUrl = videoTable.PlayUrl
	responseVideo.CoverUrl = videoTable.CoverUrl
	responseVideo.FavoriteCount = videoTable.FavoriteCount
	responseVideo.CommentCount = videoTable.CommentCount
	responseVideo.Title = videoTable.Title

	// 向数据库查询该用户是否点赞了该视频
	var err error
	go func() {
		// 先查询redis中该用户是否喜欢该视频
		favoriteRdb, err := redis.QueryFavoriteRdb(curUserId, videoTable.VideoID)
		if err != nil {
			return
		}
		if favoriteRdb == true {
			responseVideo.IsFavorite = true
			log.Println("从redis中读取用户是否喜欢该视频")
		} else {
			// 查询数据库该用户是否喜欢该视频
			_, responseVideo.IsFavorite, err = mysql.QueryFavorite(videoTable.VideoID, curUserId)
			if err != nil {
				zap.L().Error("logic/feed StuffOneVideo QueryFavorite failed", zap.Error(err))
				return
			}

			// 更新同步数据库和redis中该用户喜欢视频的数据
			if responseVideo.IsFavorite == true {
				// 需要更新redis， == false的话不需要更新

				// 数据库中查询除用户喜欢视频的id
				videosId, _ := mysql.QueryFavoriteVideos(curUserId)
				// 同步到redis中
				redis.UpdateFavoriteRdb(curUserId, videosId)
				log.Println("favoriteVideos 更新redis成功")
			}
		}
		wgStuffVideo.Done()
	}()
	// 向数据库中查询User的信息
	go func() {
		responseVideo.User, err = GetUserByIdWithCurId(videoTable.AuthorID, curUserId)
		if err != nil {
			zap.L().Error("logic/feed StuffOneVideo QueryFavorite failed", zap.Error(err))
		}
		wgStuffVideo.Done()
	}()
	wgStuffVideo.Wait()
}

// FeedSlice 对视频流进行排序相关
type FeedSlice []models.Video

func (a FeedSlice) Len() int {
	return len(a)
}
func (a FeedSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a FeedSlice) Less(i, j int) bool {
	return a[i].ID > a[j].ID
}
