package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

//声明一个全局的rdb变量
var rdbFavorite *redis.Client

// Init 初始化连接
func Init() (err error) {
	//这里别加:了，一定要注意
	rdbFavorite = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", //服务器ip地址和端口号
			viper.GetString("redis.host"),
			viper.GetInt("redis.port"),
		),
		Password: viper.GetString("redis.password"), //redis密码
		DB:       0,                                 //使用哪个数据库，0-15，选一个
		PoolSize: viper.GetInt("redis.pool_size"),   //连接池大小
	})
	return
}

func Close() {
	_ = rdbFavorite.Close()
}
