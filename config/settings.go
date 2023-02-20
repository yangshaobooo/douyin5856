package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Init 加载配置文件
func Init() (err error) {
	//方式一：直接指定配置文件路径，相对路径，绝对路径都行
	//当本地有多个不同格式的config文件时，比如config.json config.yaml
	//viper.SetConfigFile("config.yaml") 就用这个命令

	//方式二，使用文件名和文件路径配合使用
	//配置文件名不需要带后缀
	//配置文件位置可以配置多个
	viper.SetConfigName("config")   // 指定配置文件名称（不需要带后缀）
	viper.AddConfigPath("./config") //指定查找配置文件的路径（相对路径）
	//viper.SetConfigType("yaml")   // 从远程获取配置文件的时候,指定配置文件类型.对于本地的文件这个函数不生效

	err = viper.ReadInConfig() // 读取配置信息
	if err != nil {            // 读取配置信息失败
		fmt.Printf("viper.ReadInConfig failed,err:%v\n", err)
		return
	}
	viper.WatchConfig() //配置修改之后自动加载
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件需更改了...")
	})
	return
}
