package main

import (
	"fmt"
	"week07/practice/gin/controllers"
	"week07/practice/gin/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("🚀 启动学生信息管理系统...")

	// 初始化数据库 (V2 和 V3)
	controllers.InitV2DB()
	controllers.InitV3DB()

	// 初始化 Gin 引擎
	r := gin.Default()

	// 加载路由
	routes.SetupRoutes(r)

	// 启动服务，监听在 8082 端口
	if err := r.Run(":8082"); err != nil {
		fmt.Printf("服务启动失败: %v\n", err)
	}
}
