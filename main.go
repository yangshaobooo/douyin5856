package main

import (
	"douyin5856/config"
	"douyin5856/routes"
	"go.uber.org/zap"
)

func main() {
	// 1、初始化配置服务
	config.InitAllConfig()
	// 2、设置路由
	r := routes.Setup()
	// 3、启动服务
	if err := r.Run(); err != nil {
		zap.L().Error("r.run() server start failed", zap.Error(err))
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	//4、关闭开启的各种连接资源
	defer config.CloseRes()
}
