package main

import (
	"photo-album-server/dao"
	"photo-album-server/handlers"
	"photo-album-server/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	dao.InitDB()

	r := gin.Default()

	// 静态文件服务（为了能看到上传的头像）
	r.Static("/uploads", "./uploads")

	v1 := r.Group("/api/v1")
	{
		// 1. 用户模块接口 (公开) [cite: 36, 41]
		v1.POST("/register", handlers.RegisterHandler)
		v1.POST("/login", handlers.LoginHandler)

		// 2. 需要鉴权的相册模块 
		albumGroup := v1.Group("/albums")
		albumGroup.Use(middleware.AuthMiddleware())
		{
			albumGroup.POST("", handlers.CreateAlbumHandler)      // 创建 [cite: 45]
			albumGroup.GET("", handlers.GetMyAlbumsHandler)       // 我的列表 [cite: 47]
			albumGroup.GET("/public", handlers.GetPublicAlbumsHandler) // 公开广场 [cite: 51]
		}
	}

	r.Run(":8080")
}