package main

import (
	"ai-vocabulary-backend/controllers"
	"ai-vocabulary-backend/middlewares"
	
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 连接数据库
	dsn := "root:rootpassword@tcp(db:3306)/smart_vocab?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败", err)
	}

	r := gin.New() // 使用 New 而不是 Default，避免不必要的默认中间件
	r.Use(gin.Logger(), gin.Recovery())

	// 用户模块
	authRepo := &controllers.AuthHandler{DB: db}
	r.POST("/register", authRepo.Register)
	r.POST("/login", authRepo.Login)

	// 单词模块 (需 JWT 保护)
	wordRepo := &controllers.WordHandler{DB: db}
	authorized := r.Group("/")
	authorized.Use(middlewares.JWTAuth()) // JWT 中间件
	{
		authorized.GET("/query", wordRepo.QueryWord)       // 智能查询
		authorized.POST("/save", wordRepo.SaveWord)        // 手动保存
		authorized.GET("/words", wordRepo.ListWords)       // 分页列表
		authorized.DELETE("/words/:id", wordRepo.Delete)   // 软删除
	}

	r.Run(":8080")
}