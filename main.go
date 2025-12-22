package main

import (
	"book-management-system/config"
	"book-management-system/routers"
	"log"
)

func main() {
	// 加载配置
	config.LoadConfig()

	// 连接数据库
	if err := config.ConnectDatabase(); err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 设置路由
	router := routers.SetupRouter()

	// 启动服务器
	log.Printf("服务器启动在端口 %s", config.AppConfig.ServerPort)
	if err := router.Run(":" + config.AppConfig.ServerPort); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
