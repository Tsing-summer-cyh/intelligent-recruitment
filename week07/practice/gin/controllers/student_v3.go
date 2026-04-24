package controllers

import (
	"fmt"
	"log"
	"net/http"
	"week07/practice/gin/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbV3 *gorm.DB

// InitV3DB 初始化 MySQL 数据库 (需确保你的本地环境已有 MySQL 运行并创建了 practice 库)
func InitV3DB() {
	// 请替换为你自己的 MySQL 用户名和密码
	dsn := "root:123456@tcp(127.0.0.1:3306)/practice?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	dbV3, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("[V3] MySQL 数据库连接失败: %v\n(如未配置 MySQL，V3 接口将不可用)", err)
		return
	}

	dbV3.AutoMigrate(&models.Student{})
	fmt.Println("[V3] MySQL 数据库连接且自动迁移成功")
}

func CreateStudentV3(c *gin.Context) {
	var req models.Student
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dbV3.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	fmt.Printf("[V3] 成功创建学生: %+v\n", req)
	c.JSON(http.StatusCreated, gin.H{"message": "创建成功", "data": req})
}

func GetStudentsV3(c *gin.Context) {
	var students []models.Student
	dbV3.Find(&students)
	c.JSON(http.StatusOK, students)
}

func DeleteStudentV3(c *gin.Context) {
	id := c.Param("id")
	if err := dbV3.Delete(&models.Student{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	fmt.Printf("[V3] 成功删除学生 ID: %s\n", id)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
