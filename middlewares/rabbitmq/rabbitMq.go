package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

const MqAddress = "amqp://admin:admin@localhost:5672/"

type RabbitMQ struct {
	conn      *amqp.Connection
	mqAddress string
}

var Rmq *RabbitMQ

// InitRabbitMQ 初始化RabbitMQ的连接和通道
func InitRabbitMQ() {
	Rmq = &RabbitMQ{
		mqAddress: MqAddress,
	}
	dial, err := amqp.Dial(Rmq.mqAddress)
	Rmq.failOnErr(err, "MQ创建连接失败")
	Rmq.conn = dial
}

// 连接出错时，输出错误信息。
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s\n", err, message)
		panic(fmt.Sprintf("%s:%s\n", err, message))
	}
}

// 关闭mq通道和mq的连接。
func (r *RabbitMQ) destroy() {
	r.conn.Close()
}
