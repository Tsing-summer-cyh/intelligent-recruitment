package db

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// 注意替换为你自己的 MySQL 密码
	dsn := "root:123456@tcp(127.0.0.1:3306)/recruitment_db?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接 MySQL 失败: %v", err)
	}

	fmt.Println("🚀 MySQL 数据库连接成功！")

	// 1. 自动同步表结构 (AutoMigrate)
	err = DB.AutoMigrate(&User{}, &Job{}, &CandidateProfile{}, &JobApplication{}, &AiChatHistory{})
	if err != nil {
		log.Fatalf("自动建表失败: %v", err)
	}
	fmt.Println("📦 数据表结构已自动同步！")

	// 2. 注入默认测试账号 (如果不存在的话)
	var count int64

	// HR 测试账号
	DB.Model(&User{}).Where("username = ?", "test_hr").Count(&count)
	if count == 0 {
		testUser := User{
			Username:     "test_hr",
			PasswordHash: "123456", // 演示用明文
			Role:         "hr",
		}
		DB.Create(&testUser)
		fmt.Println("👤 已自动创建测试 HR 账号: test_hr / 123456")
	}

	// 候选人测试账号
	DB.Model(&User{}).Where("username = ?", "test_candidate").Count(&count)
	if count == 0 {
		candidateUser := User{
			Username:     "test_candidate",
			PasswordHash: "123456", // 演示用明文
			Role:         "candidate",
		}
		DB.Create(&candidateUser)
		fmt.Println("👤 已自动创建测试候选人账号: test_candidate / 123456")
	}
}