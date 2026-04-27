package controllers

import (
	"ai-vocabulary-backend/models"
	"ai-vocabulary-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WordHandler struct {
	DB *gorm.DB
}

// QueryWord 智能查询单词 (作业核心逻辑)
func (h *WordHandler) QueryWord(c *gin.Context) {
	// 从 JWT 中间件获取当前登录用户的 ID
	userID, _ := c.Get("userID") 
	word := c.Query("word")
	aiProvider := c.Query("ai_provider")

	if word == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "单词不能为空"})
		return
	}
	if aiProvider == "" {
		aiProvider = "DeepSeek" // 默认备用模型
	}

	// 1. 检查数据库中当前用户是否已保存过该单词
	var existingWord models.Word
	if err := h.DB.Where("user_id = ? AND word = ?", userID, word).First(&existingWord).Error; err == nil {
		// 已保存，直接从数据库读取并返回
		c.JSON(http.StatusOK, gin.H{
			"source": "database",
			"data":   existingWord,
		})
		return
	}

	// 2. 未保存，调用 AI 接口生成释义和例句
	meaning, sentences, err := services.FetchAIResult(word, aiProvider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI 调用失败", "details": err.Error()})
		return
	}

	// 3. 直接将 AI 结果返回给前端展示，此时【不在后端进行数据库保存】(严格遵循作业要求)
	c.JSON(http.StatusOK, gin.H{
		"source": "ai",
		"data": gin.H{
			"word":        word,
			"meaning":     meaning,
			"sentences":   sentences, // 这是一个数组
			"ai_provider": aiProvider,
		},
	})
}

// SaveWord 手动保存单词
func (h *WordHandler) SaveWord(c *gin.Context) {
	userID, _ := c.Get("userID")

	// 接收前端提交的完整数据
	var req struct {
		Word       string   `json:"word" binding:"required"`
		Meaning    string   `json:"meaning" binding:"required"`
		Sentences  []string `json:"sentences" binding:"required"`
		AIProvider string   `json:"ai_provider" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误，请检查提交格式"})
		return
	}

	// 组装数据并存入 MySQL (与 UserID 绑定)
	wordRecord := models.Word{
		UserID:     userID.(uint),
		Word:       req.Word,
		Meaning:    req.Meaning,
		Sentences:  models.SentenceList(req.Sentences), // 转换为自定义的 JSON 类型
		AIProvider: req.AIProvider,
	}

	if err := h.DB.Create(&wordRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存单词失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "保存成功", "data": wordRecord})
}

// ListWords 获取单词列表 (要求支持分页)
func (h *WordHandler) ListWords(c *gin.Context) {
	userID, _ := c.Get("userID")

	// 获取分页参数，默认第 1 页，每页 10 条
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	var words []models.Word
	var total int64

	// 查询总记录数
	h.DB.Model(&models.Word{}).Where("user_id = ?", userID).Count(&total)

	// 查询分页数据
	if err := h.DB.Where("user_id = ?", userID).Offset(offset).Limit(pageSize).Order("created_at desc").Find(&words).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取单词列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"data":      words,
	})
}

// Delete 软删除单词记录
func (h *WordHandler) Delete(c *gin.Context) {
	userID, _ := c.Get("userID")
	wordID := c.Param("id")

	// GORM 默认开启了软删除 (因为模型里包含 DeletedAt 字段)
	if err := h.DB.Where("id = ? AND user_id = ?", wordID, userID).Delete(&models.Word{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}