package config

import (
	"douyin5856/dao/mysql"
	"douyin5856/dao/redis"
	"douyin5856/logger"
	"douyin5856/middlewares/snowflake1"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitAllConfig() {
	//1、加载配置  再setting中的Init
	if err := Init(); err != nil {
		fmt.Printf("Init settings failed,err:%v\n", err)
		return
	}
	//2、初始化日志
	if err := logger.Init(); err != nil {
		fmt.Printf("Init settings failed,err:%v\n", err)
		return
	}
	//3、初始化连接Mysql连接
	if err := mysql.Init(); err != nil {
		fmt.Printf("Init mysql failed,err:%v\n", err)
		return
	}

	//4、初始化Redis连接
	if err := redis.Init(); err != nil {
		fmt.Printf("Init redis failed,err:%v\n", err)
		return
	}

	// 5、初始化雪花算法，用于获取随机id
	if err := snowflake1.Init(viper.GetString("app.start_time"), viper.GetInt64("app.machine_id")); err != nil {
		fmt.Printf("init snowflake1 failed, err:%v\n", err)
		return
	}
	zap.L().Info("InitAllConfig(settings/logger/mysql/redis/snowflake1) success")
}

func CloseRes() {
	defer zap.L().Sync()
	defer mysql.Close()
	defer redis.Close()
}
