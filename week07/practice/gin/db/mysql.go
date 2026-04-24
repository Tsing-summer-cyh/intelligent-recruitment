package db

import (
	"log"
	"week07/practice/gin/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MySQL *gorm.DB

func InitMySQL() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/practice?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	MySQL, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("MySQL 连接失败: %v", err)
	}

	MySQL.AutoMigrate(&models.Student{})
	log.Println("✅ MySQL 初始化成功")
}
