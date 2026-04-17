package handlers

import (
	"net/http"
	"photo-album-server/dao"
	"photo-album-server/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateAlbumHandler 创建相册 [cite: 45-46]
func CreateAlbumHandler(c *gin.Context) {
	userID := c.MustGet("userID").(int) // 从中间件获取

	var album models.Album
	if err := c.ShouldBindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的输入数据"})
		return
	}

	_, err := dao.DB.Exec("INSERT INTO albums (user_id, name, description, is_public) VALUES (?, ?, ?, ?)",
		userID, album.Name, album.Description, album.IsPublic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建相册失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "相册创建成功"})
}

// GetMyAlbumsHandler 获取“我的”相册列表 [cite: 47-50]
func GetMyAlbumsHandler(c *gin.Context) {
	userID := c.MustGet("userID").(int)
	
	// 分页处理 [cite: 50]
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	offset := (page - 1) * size

	rows, err := dao.DB.Query("SELECT id, name, description, is_public, created_at FROM albums WHERE user_id = ? LIMIT ? OFFSET ?",
		userID, size, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	defer rows.Close()

	var albums []models.Album
	for rows.Next() {
		var a models.Album
		rows.Scan(&a.ID, &a.Name, &a.Description, &a.IsPublic, &a.CreatedAt)
		albums = append(albums, a)
	}

	c.JSON(http.StatusOK, albums)
}

// GetPublicAlbumsHandler 获取“公开广场”相册 [cite: 51-55]
func GetPublicAlbumsHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	offset := (page - 1) * size

	// 必须包含创建者的 username 和 avatar_url [cite: 54]
	query := `
		SELECT a.id, a.name, a.description, a.created_at, u.username, u.avatar_url
		FROM albums a
		JOIN users u ON a.user_id = u.id
		WHERE a.is_public = 1
		LIMIT ? OFFSET ?`
	
	rows, err := dao.DB.Query(query, size, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	defer rows.Close()

	var albums []models.Album
	for rows.Next() {
		var a models.Album
		rows.Scan(&a.ID, &a.Name, &a.Description, &a.CreatedAt, &a.Username, &a.AvatarURL)
		albums = append(albums, a)
	}

	c.JSON(http.StatusOK, albums)
}