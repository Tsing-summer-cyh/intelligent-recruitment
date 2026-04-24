package v4

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"week07/practice/gin/db"
	"week07/practice/gin/models"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	CacheExpiration = 5 * time.Minute // 缓存过期时间
	AllStudentsKey  = "students:all"  // 所有学生的缓存 Key
)

// getStudentKey 生成单个学生的缓存 Key
func getStudentKey(id string) string {
	return "student:" + id
}

// CreateStudent 创建学生 (写操作：清空列表缓存)
func CreateStudent(c *gin.Context) {
	var req models.Student
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	if err := db.MySQL.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库写入失败"})
		return
	}

	// 缓存一致性：新增了数据，原来的“所有学生”缓存就失效了，必须删除
	db.Redis.Del(db.Ctx, AllStudentsKey)

	log.Printf("[V4] 成功创建学生并清理列表缓存: %+v\n", req)
	c.JSON(http.StatusCreated, gin.H{"message": "创建成功", "data": req})
}

// GetStudents 获取所有学生 (读操作：先读缓存，未命中读库并写入缓存)
func GetStudents(c *gin.Context) {
	// 1. 尝试从 Redis 获取
	cacheData, err := db.Redis.Get(db.Ctx, AllStudentsKey).Result()
	if err == nil {
		// 缓存命中
		var students []models.Student
		json.Unmarshal([]byte(cacheData), &students)
		log.Println("[V4] 从 Redis 缓存中获取了所有学生数据 🎯")
		c.JSON(http.StatusOK, students)
		return
	} else if err != redis.Nil {
		log.Printf("Redis 错误: %v\n", err)
	}

	// 2. 缓存未命中，从 MySQL 获取
	var students []models.Student
	db.MySQL.Find(&students)

	// 3. 将结果序列化并存入 Redis，设置 5 分钟过期
	jsonBytes, _ := json.Marshal(students)
	db.Redis.Set(db.Ctx, AllStudentsKey, jsonBytes, CacheExpiration)

	log.Println("[V4] 从 MySQL 获取数据，并写入 Redis 缓存 🗄️")
	c.JSON(http.StatusOK, students)
}

// GetStudent 获取单个学生
func GetStudent(c *gin.Context) {
	id := c.Param("id")
	cacheKey := getStudentKey(id)

	// 1. 查缓存
	cacheData, err := db.Redis.Get(db.Ctx, cacheKey).Result()
	if err == nil {
		var student models.Student
		json.Unmarshal([]byte(cacheData), &student)
		log.Printf("[V4] 从 Redis 获取学生(ID:%s) 🎯\n", id)
		c.JSON(http.StatusOK, student)
		return
	}

	// 2. 查库
	var student models.Student
	if err := db.MySQL.First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到该学生"})
		return
	}

	// 3. 写入缓存
	jsonBytes, _ := json.Marshal(student)
	db.Redis.Set(db.Ctx, cacheKey, jsonBytes, CacheExpiration)

	log.Printf("[V4] 从 MySQL 获取学生(ID:%s)并缓存 🗄️\n", id)
	c.JSON(http.StatusOK, student)
}

// UpdateStudent 更新学生 (写操作：双删缓存)
func UpdateStudent(c *gin.Context) {
	id := c.Param("id")
	var req models.Student
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 检查是否存在
	var student models.Student
	if err := db.MySQL.First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到该学生"})
		return
	}

	// 更新 (排除 ID 的更新)
	db.MySQL.Model(&student).Updates(req)

	// 缓存一致性：删除单个学生缓存 + 列表缓存
	db.Redis.Del(db.Ctx, getStudentKey(id), AllStudentsKey)

	log.Printf("[V4] 成功更新学生(ID:%s)并清理相关缓存\n", id)
	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "data": student})
}

// DeleteStudent 删除学生 (写操作：双删缓存)
func DeleteStudent(c *gin.Context) {
	id := c.Param("id")

	if err := db.MySQL.Delete(&models.Student{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	// 缓存一致性：删除相关缓存
	db.Redis.Del(db.Ctx, getStudentKey(id), AllStudentsKey)

	log.Printf("[V4] 成功删除学生(ID:%s)并清理相关缓存\n", id)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
