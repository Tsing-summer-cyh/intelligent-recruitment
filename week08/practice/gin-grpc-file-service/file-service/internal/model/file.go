// file-service/internal/model/file.go
package model

import "time"

// FileRecord 对应数据库中的 file_records 表
type FileRecord struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"`
	OriginalName string    `gorm:"type:varchar(255);not null"` // 原始文件名 [cite: 84, 88]
	StoredName   string    `gorm:"type:varchar(255);not null;unique"` // 哈希后的文件名 [cite: 80, 89]
	Size         int64     `gorm:"not null"`                   // 文件大小 [cite: 90]
	MimeType     string    `gorm:"type:varchar(100)"`          // 文件类型 [cite: 91]
	Path         string    `gorm:"type:varchar(255);not null"` // 在 uploads 目录下的相对路径 [cite: 92]
	CreatedAt    time.Time // 记录创建时间
}