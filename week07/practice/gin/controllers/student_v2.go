package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"week07/practice/gin/models"

	"github.com/gin-gonic/gin"
	// 1. 换成纯 Go 版本的驱动
	_ "modernc.org/sqlite"
)

var dbV2 *sql.DB

// InitV2DB 初始化 SQLite 数据库
func InitV2DB() {
	var err error
	// 2. 这里的驱动名称改为 "sqlite"
	dbV2, err = sql.Open("sqlite", "./students.db")
	if err != nil {
		log.Fatalf("[V2] SQLite 数据库连接失败: %v", err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		age INTEGER NOT NULL,
		grade TEXT NOT NULL
	);`
	_, err = dbV2.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("[V2] 创建表失败: %v", err)
	}
	fmt.Println("[V2] SQLite 数据库初始化成功")
}

func CreateStudentV2(c *gin.Context) {
	var req models.Student
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := dbV2.Exec("INSERT INTO students (name, age, grade) VALUES (?, ?, ?)", req.Name, req.Age, req.Grade)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "插入数据失败"})
		return
	}
	id, _ := res.LastInsertId()
	req.ID = int(id)

	fmt.Printf("[V2] 成功创建学生: %+v\n", req)
	c.JSON(http.StatusCreated, gin.H{"message": "创建成功", "data": req})
}

func GetStudentsV2(c *gin.Context) {
	rows, err := dbV2.Query("SELECT id, name, age, grade FROM students")
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

// 省略了 V2 的 GetOne, Update, Delete 具体实现以保持精简，
// 逻辑与 GetStudentsV2 类似，使用 dbV2.QueryRow 和 dbV2.Exec 即可。
