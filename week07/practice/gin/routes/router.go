package routes

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine) {
	// 注册各个版本的路由
	RegisterV1Routes(r)
	RegisterV2Routes(r)
	RegisterV3Routes(r)
	RegisterV4Routes(r) // 这是上一轮给你的 V4 路由
}
