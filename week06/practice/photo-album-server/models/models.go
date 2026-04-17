package models

import "time"

// User 用户表结构 [cite: 10-16]
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // 密码在 JSON 中隐藏
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
}

// Album 相册表结构 [cite: 17-25]
type Album struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsPublic    int       `json:"is_public"` // 0: 私有, 1: 公有 [cite: 22]
	CreatedAt   time.Time `json:"created_at"`
	// 扩展字段：用于公开广场展示创建者信息 [cite: 54]
	Username  string `json:"username,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

// Photo 照片表结构 [cite: 26-33]
type Photo struct {
	ID        int       `json:"id"`
	AlbumID   int       `json:"album_id"`
	FilePath  string    `json:"file_path"`
	FileSize  int64     `json:"file_size"`
	CreatedAt time.Time `json:"created_at"`
}