package rabbitmq

import (
	"douyin5856/dao/mysql"
	"fmt"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"log"
	"strconv"
	"strings"
)

type LikeMQ struct {
	RabbitMQ
	channel   *amqp.Channel
	queueName string
	exchange  string
	key       string
}

// NewLikeRabbitMQ 获取likeMQ的对应队列。
func NewLikeRabbitMQ(queueName string) *LikeMQ {
	likeMQ := &LikeMQ{
		RabbitMQ:  *Rmq,
		queueName: queueName,
	}
	cha, err := likeMQ.conn.Channel()
	likeMQ.channel = cha
	Rmq.failOnErr(err, "获取通道失败")
	return likeMQ
}

// Publish like操作的发布配置。
func (l *LikeMQ) Publish(message string) {

	_, err := l.channel.QueueDeclare(
		l.queueName,
		//是否持久化
		false,
		//是否为自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		//额外属性
		nil,
	)
	if err != nil {
		panic(err)
	}

	err1 := l.channel.Publish(
		l.exchange,
		l.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err1 != nil {
		panic(err)
	}

}

// Consumer like关系的消费逻辑。
func (l *LikeMQ) Consumer() {

	_, err := l.channel.QueueDeclare(l.queueName, false, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	//2、接收消息
	messages, err1 := l.channel.Consume(
		l.queueName,
		//用来区分多个消费者
		"",
		//是否自动应答
		true,
		//是否具有排他性
		false,
		//如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		//消息队列是否阻塞
		false,
		nil,
	)
	if err1 != nil {
		panic(err1)
	}

	forever := make(chan bool)
	switch l.queueName {
	case "like_add":
		//点赞消费队列
		go l.consumerLikeAdd(messages)
	case "like_del":
		//取消赞消费队列
		go l.consumerLikeDel(messages)

	}

	log.Printf("[*] Waiting for messagees,To exit press CTRL+C")

	<-forever

}

//consumerLikeAdd 赞关系添加的消费方式。
func (l *LikeMQ) consumerLikeAdd(messages <-chan amqp.Delivery) {
	for d := range messages {
		// 参数解析
		params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
		userId, _ := strconv.ParseInt(params[0], 10, 64)
		videoId, _ := strconv.ParseInt(params[1], 10, 64)
		maxAttempts := viper.GetInt("mysql.max_attempts")
		// 最多尝试操作数据库的次数
		for i := 0; i < maxAttempts; i++ {
			flag := false //默认没有错误
			// 1、判断该条点赞信息是否存在过
			haven, isFavorite, err := mysql.QueryFavorite(videoId, userId)
			if err != nil {
				flag = true // 出现问题
				zap.L().Error("likeMq mysql.QueryFavorite failed", zap.Error(err))
			}
			// 2、如果该条数据存在过
			if haven {
				isFavorite = !isFavorite
				// 3、修改用户点赞表中的数据
				if err = mysql.AlterFavorite(userId, videoId, isFavorite); err != nil {
					flag = true
					zap.L().Error("likeMq mysql.AlterFavorite failed,", zap.Error(err))
				}
				// 4、修改视频的点赞数量，+1
				if err = mysql.AddVideoFavoriteCount(videoId, 1); err != nil {
					flag = true
					zap.L().Error("likeMq mysql.AddFavoriteCount failed,", zap.Error(err))
				}
				// 用户点赞视频数量的字段+1
				mysql.AddUserFavoriteCount(userId, 1)
				// 找到该视频的用户
				videoUserId, _ := mysql.QueryUserIdByVideoId(videoId)
				// 视频用户获赞总数+1
				mysql.AddAllFavoriteCount(videoUserId, 1)

				log.Println("点赞点赞")
			} else {
				// 该条数据没有存在过
				isFavorite = true
				// 5、插入新数据
				if err = mysql.InsertFavorite(userId, videoId, isFavorite); err != nil {
					flag = true
					zap.L().Error("likeMq mysql.InsertFavorite failed,", zap.Error(err))
				}
				// 6、 修改该视频点赞数量，+1
				if err = mysql.AddVideoFavoriteCount(videoId, 1); err != nil {
					flag = true
					zap.L().Error("likeMq mysql.AddFavoriteCount failed,", zap.Error(err))
				}
				// 用户点赞视频数量的字段+1
				mysql.AddUserFavoriteCount(userId, 1)
				// 找到该视频的用户
				videoUserId, _ := mysql.QueryUserIdByVideoId(videoId)
				// 视频用户获赞总数+1
				mysql.AddAllFavoriteCount(videoUserId, 1)

				log.Println("点赞点赞")
			}
			// 7、一遍流程下来正常执行了，那就打断结束，不再尝试
			if flag == false {
				break
			}
		}
	}
}

//consumerLikeDel 赞关系删除的消费方式。
func (l *LikeMQ) consumerLikeDel(messages <-chan amqp.Delivery) {
	for d := range messages {
		// 参数解析。
		params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
		userId, _ := strconv.ParseInt(params[0], 10, 64)
		videoId, _ := strconv.ParseInt(params[1], 10, 64)
		//最多尝试操作数据库的次数
		maxAttempts := viper.GetInt("mysql.max_attempts")
		for i := 0; i < maxAttempts; i++ {
			flag := false //默认没有错误
			// 1、判断该条点赞信息是否存在过
			haven, isFavorite, err := mysql.QueryFavorite(videoId, userId)
			if err != nil {
				flag = true // 出现问题
				zap.L().Error("likeMq mysql.QueryFavorite failed", zap.Error(err))
			}
			// 2、如果该条数据存在过
			if haven {
				isFavorite = !isFavorite
				// 3、修改用户点赞表中的数据
				if err = mysql.AlterFavorite(userId, videoId, isFavorite); err != nil {
					flag = true
					zap.L().Error("likeMq mysql.AlterFavorite failed,", zap.Error(err))
				}
				// 4、修改视频的点赞数量，-1
				if err = mysql.AddVideoFavoriteCount(videoId, -1); err != nil {
					flag = true
					zap.L().Error("likeMq mysql.AddFavoriteCount failed,", zap.Error(err))
				}
				// 用户点赞视频数量的字段 -1
				mysql.AddUserFavoriteCount(userId, -1)
				// 找到该视频的用户
				videoUserId, _ := mysql.QueryUserIdByVideoId(videoId)
				// 视频用户获赞总数-1
				mysql.AddAllFavoriteCount(videoUserId, -1)
				log.Println("取消点赞")
			} else {
				// 该条数据没有存在过
				isFavorite = false
				// 5、插入新数据
				if err = mysql.InsertFavorite(userId, videoId, isFavorite); err != nil {
					flag = true
					zap.L().Error("likeMq mysql.InsertFavorite failed,", zap.Error(err))
				}
				// 6、 修改该视频点赞数量，-1
				if err = mysql.AddVideoFavoriteCount(videoId, -1); err != nil {
					flag = true
					zap.L().Error("likeMq mysql.AddFavoriteCount failed,", zap.Error(err))
				}
				// 用户点赞视频数量的字段+1
				mysql.AddUserFavoriteCount(userId, -1)
				// 找到该视频的用户
				videoUserId, _ := mysql.QueryUserIdByVideoId(videoId)
				// 视频用户获赞总数+1
				mysql.AddAllFavoriteCount(videoUserId, -1)
				log.Println("取消点赞")
			}
			// 7、一遍流程下来正常执行了，那就打断结束，不再尝试
			if flag == false {
				break
			}
		}
	}
}

var RmqLikeAdd *LikeMQ
var RmqLikeDel *LikeMQ

// InitLikeRabbitMQ 初始化rabbitMQ连接。
func InitLikeRabbitMQ() {
	RmqLikeAdd = NewLikeRabbitMQ("like_add")
	go RmqLikeAdd.Consumer()

	RmqLikeDel = NewLikeRabbitMQ("like_del")
	go RmqLikeDel.Consumer()
}
