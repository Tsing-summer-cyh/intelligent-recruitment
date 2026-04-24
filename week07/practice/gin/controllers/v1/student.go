package v1

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"week07/practice/gin/models"

	"github.com/gin-gonic/gin"
)

var (
	studentsV1 = make([]models.Student, 0)
	idCounter  = 1
	mutex      sync.Mutex // 加锁防止并发冲突
)

func CreateStudent(c *gin.Context) {
	var req models.Student
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据绑定失败: " + err.Error()})
		return
	}

	mutex.Lock()
	req.ID = idCounter
	idCounter++
	studentsV1 = append(studentsV1, req)
	mutex.Unlock()

	log.Printf("[V1] 成功创建学生: %+v\n", req)
	c.JSON(http.StatusCreated, gin.H{"message": "创建成功", "data": req})
}

func GetStudents(c *gin.Context) {
	log.Printf("[V1] 获取所有学生信息，当前数量: %d\n", len(studentsV1))
	c.JSON(http.StatusOK, studentsV1)
}

func GetStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID 格式"})
		return
	}

	for _, s := range studentsV1 {
		if s.ID == id {
			log.Printf("[V1] 获取单个学生成功: %+v\n", s)
			c.JSON(http.StatusOK, s)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "未找到该学生"})
}

func UpdateStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var req models.Student
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据绑定失败"})
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for i, s := range studentsV1 {
		if s.ID == id {
			req.ID = id // 保持 ID 不变
			studentsV1[i] = req
			log.Printf("[V1] 更新学生成功: %+v\n", req)
			c.JSON(http.StatusOK, gin.H{"message": "更新成功", "data": req})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "未找到该学生"})
}

func DeleteStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	mutex.Lock()
	defer mutex.Unlock()

	for i, s := range studentsV1 {
		if s.ID == id {
			studentsV1 = append(studentsV1[:i], studentsV1[i+1:]...)
			log.Printf("[V1] 删除学生成功，ID: %d\n", id)
			c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "未找到该学生"})
}
