package utils

import (
	
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 签名密钥 (实际项目中应写在配置文件里，这里为演示直接写死)
var jwtSecret = []byte("whut_smart_recruitment_secret_key_2026")

// CustomClaims 自定义载荷，包含我们需要存放在 Token 里的用户信息
type CustomClaims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT 令牌
func GenerateToken(userID int64, role string) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token 24小时有效
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// JWTAuthMiddleware JWT 鉴权中间件
func JWTAuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头获取 Authorization 字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未携带 Token，无权访问"})
			c.Abort()
			return
		}

		// 2. 按空格分割 (规范是 Bearer <token>)
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Token 格式错误"})
			c.Abort()
			return
		}

		// 3. 解析 Token
		tokenString := parts[1]
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Token 无效或已过期"})
			c.Abort()
			return
		}

		// 4. 获取用户信息并校验角色
		claims, ok := token.Claims.(*CustomClaims)
		if ok {
			// 如果 requiredRole 不为空，且当前角色不匹配，则拒绝访问
			if requiredRole != "" && claims.Role != requiredRole {
				c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "权限不足：需要 " + requiredRole + " 角色"})
				c.Abort()
				return
			}
			// 将用户 ID 和 Role 存入上下文，方便后续业务逻辑使用
			c.Set("userID", claims.UserID)
			c.Set("role", claims.Role)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Token 解析失败"})
			c.Abort()
		}
	}
}