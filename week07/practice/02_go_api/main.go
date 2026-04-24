package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// 设置为 release 模式，让终端输出清爽一点
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 定义 /ping 路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 监听并在 0.0.0.0:8081 上启动服务
	r.Run(":8081")
}
