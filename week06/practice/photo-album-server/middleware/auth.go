package middleware

import (
	"net/http"
	"photo-album-server/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从 Header 获取 Authorization [cite: 43]
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "请提供有效的鉴权 Token"})
			c.Abort() // 拦截后续处理
			return
		}

		// 2. 解析 Token
		tokenString := authHeader[7:] // 去掉 "Bearer " 前缀
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token 无效或已过期"})
			c.Abort()
			return
		}

		// 3. 将解析出的 userID 存入上下文，方便后续接口使用
		c.Set("userID", claims.UserID)
		c.Next()
	}
}