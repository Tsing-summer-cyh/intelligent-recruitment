// file-web/cmd/main.go
package main

import (
	"log"

	"file-web/internal/grpcclient"
	"file-web/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 初始化 gRPC 客户端，连接到 file-service
	grpcclient.InitClient()

	// 2. 创建 Gin 引擎
	r := gin.Default()

	// 3. 注册路由组 [cite: 75, 95, 108]
	api := r.Group("/api/files")
	{
		api.POST("/uploads", handler.UploadFiles)
		api.GET("", handler.GetFiles)
		api.GET("/download/:id", handler.DownloadFile)
	}

	log.Println("Web 服务已启动，正在监听端口 :8080 ...")
	
	// 4. 启动 HTTP 服务
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("启动 Web 服务失败: %v", err)
	}
}