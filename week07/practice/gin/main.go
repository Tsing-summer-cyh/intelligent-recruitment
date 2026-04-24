package main

import (
	"log"
	"week07/practice/gin/db"
	"week07/practice/gin/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("🚀 启动学生信息管理系统 (V4版)...")

	// 1. 初始化数据库及缓存组件
	// db.InitSQLite() // V2 使用
	db.InitMySQL() // V3/V4 使用
	db.InitRedis() // V4 使用

	// 2. 初始化 Gin 引擎
	r := gin.Default()

	// 3. 注册所有路由
	routes.SetupRoutes(r)

	// 4. 启动服务
	if err := r.Run(":8082"); err != nil {
		log.Fatalf("服务启动失败: %v\n", err)
	}
}
