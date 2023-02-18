package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

//声明一个全局的sqlx数据库
var db *sqlx.DB

func Init() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"),
	)
	//Connect包含了open和ping
	//也可以使用MustConnect，连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed", zap.Error(err))
		return
	}
	//设置最大连接数和空闲数
	db.SetMaxOpenConns(viper.GetInt("mysql.max_open_connections"))
	db.SetMaxIdleConns(viper.GetInt("mysql.max_idle_connections"))
	return
}

func Close() {
	_ = db.Close()
}
