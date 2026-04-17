package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"photo-album-server/dao"
	"photo-album-server/models"
	"photo-album-server/utils"
	"time"
	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	file, _ := c.FormFile("avatar")

	if file != nil {
		ext := filepath.Ext(file.Filename)
		if ext != ".jpg" && ext != ".png" && ext != ".gif" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "仅支持 jpg/png/gif 格式"})
			return
		}
		// 限制 2MB [cite: 40]
		if file.Size > 2*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "文件超过 2MB"})
			return
		}

		fileName := fmt.Sprintf("%d_%d%s", time.Now().UnixNano(), time.Now().Unix(), ext)
		savePath := filepath.Join("uploads/avatars", fileName)
		c.SaveUploadedFile(file, savePath)

		_, err := dao.DB.Exec("INSERT INTO users (username, password, avatar_url) VALUES (?, ?, ?)", 
			username, password, savePath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "注册失败，用户名可能已存在"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

func LoginHandler(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var user models.User
	err := dao.DB.QueryRow("SELECT id FROM users WHERE username = ? AND password = ?", 
		loginData.Username, loginData.Password).Scan(&user.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	token, _ := utils.GenerateToken(user.ID)
	c.JSON(http.StatusOK, gin.H{"token": token})
}