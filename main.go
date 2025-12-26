package main

import (
	"book-management-system/config"
	"book-management-system/routers"
	"fmt"
	"log"
)

// @title          图书管理系统 API
// @version        1.0
// @description    这是一个图书管理系统的REST API文档
// @description    包含用户认证、图书管理、借阅管理等功能
// @contact.name   API支持
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 输入"Bearer {token}"，token在登录后获得

// @schemes http

func main() {
	// 加载配置
	config.LoadConfig()

	// 调试：打印配置信息
	fmt.Printf("JWT Secret长度: %d\n", len(config.AppConfig.JWTSecret))
	fmt.Printf("JWT Expire: %v\n", config.AppConfig.JWTExpire)

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
