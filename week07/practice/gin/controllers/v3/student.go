package v3

import (
	"log"
	"net/http"
	"week07/practice/gin/db"
	"week07/practice/gin/models"

	"github.com/gin-gonic/gin"
)

func CreateStudent(c *gin.Context) {
	var req models.Student
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.MySQL.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}
	log.Printf("[V3] 成功创建学生: %+v\n", req)
	c.JSON(http.StatusCreated, gin.H{"message": "创建成功", "data": req})
}

func GetStudents(c *gin.Context) {
	var students []models.Student
	db.MySQL.Find(&students)
	c.JSON(http.StatusOK, students)
}

func GetStudent(c *gin.Context) {
	id := c.Param("id")
	var student models.Student
	if err := db.MySQL.First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到该学生"})
		return
	}
	c.JSON(http.StatusOK, student)
}

func UpdateStudent(c *gin.Context) {
	id := c.Param("id")
	var req models.Student
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据绑定失败"})
		return
	}

	var student models.Student
	if err := db.MySQL.First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到该学生"})
		return
	}

	db.MySQL.Model(&student).Updates(req)
	log.Printf("[V3] 成功更新学生 ID: %s\n", id)
	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "data": student})
}

func DeleteStudent(c *gin.Context) {
	id := c.Param("id")
	if err := db.MySQL.Delete(&models.Student{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	log.Printf("[V3] 成功删除学生 ID: %s\n", id)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
