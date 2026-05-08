// file-service/internal/database/db.go
package database

import (
	"log"
	"os"
	"path/filepath"

	"file-service/internal/model"

	// 替换成了这个纯 Go 驱动 👇
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// 确保 data 目录存在
	err := os.MkdirAll("data", os.ModePerm)
	if err != nil {
		log.Fatalf("创建 data 目录失败: %v", err)
	}

	dbPath := filepath.Join("data", "file.db")

	// 连接 SQLite 数据库 (这里的 sqlite.Open 调用的已经是新驱动了)
	database, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接 SQLite 数据库失败: %v", err)
	}

	// 自动迁移表结构（创建 file_records 表）
	err = database.AutoMigrate(&model.FileRecord{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	DB = database
	log.Println("SQLite 数据库初始化成功，表结构已就绪！")
}
