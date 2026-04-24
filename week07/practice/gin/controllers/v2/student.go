package v2

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

	res, err := db.SQLite.Exec("INSERT INTO students (name, age, grade) VALUES (?, ?, ?)", req.Name, req.Age, req.Grade)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "插入数据失败"})
		return
	}
	id, _ := res.LastInsertId()
	req.ID = int(id)

	log.Printf("[V2] 成功创建学生: %+v\n", req)
	c.JSON(http.StatusCreated, gin.H{"message": "创建成功", "data": req})
}

func GetStudents(c *gin.Context) {
	rows, err := db.SQLite.Query("SELECT id, name, age, grade FROM students")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var s models.Student
		rows.Scan(&s.ID, &s.Name, &s.Age, &s.Grade)
		students = append(students, s)
	}

	if students == nil {
		students = []models.Student{}
	}
	c.JSON(http.StatusOK, students)
}

func GetStudent(c *gin.Context) {
	id := c.Param("id")
	var s models.Student

	err := db.SQLite.QueryRow("SELECT id, name, age, grade FROM students WHERE id = ?", id).Scan(&s.ID, &s.Name, &s.Age, &s.Grade)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到该学生"})
		return
	}
	c.JSON(http.StatusOK, s)
}

func UpdateStudent(c *gin.Context) {
	id := c.Param("id")
	var req models.Student
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据绑定失败"})
		return
	}

	_, err := db.SQLite.Exec("UPDATE students SET name=?, age=?, grade=? WHERE id=?", req.Name, req.Age, req.Grade, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	log.Printf("[V2] 更新学生成功 ID: %s\n", id)
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func DeleteStudent(c *gin.Context) {
	id := c.Param("id")
	_, err := db.SQLite.Exec("DELETE FROM students WHERE id=?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	log.Printf("[V2] 成功删除学生 ID: %s\n", id)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
