package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your_secret_key") // 实际开发建议放入 .env

// JWTAuth 拦截器
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "请求头中缺少授权 Token 或格式错误"})
			c.Abort() // 阻止后续请求
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的或已过期的 Token"})
			c.Abort()
			return
		}

		// 解析成功，将 UserID 存入上下文供后续业务使用
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("userID", uint(claims["user_id"].(float64)))
		}
		c.Next()
	}
}