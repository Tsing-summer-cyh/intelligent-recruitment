package routes

import (
	"week07/practice/gin/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 配置所有路由
func SetupRoutes(r *gin.Engine) {
	// V1: 内存存储 (同时映射根目录 /students 以满足基本需求)
	v1 := r.Group("/students")
	{
		v1.POST("", controllers.CreateStudentV1)
		v1.GET("", controllers.GetStudentsV1)
		v1.GET("/:id", controllers.GetStudentV1)
		v1.PUT("/:id", controllers.UpdateStudentV1)
		v1.DELETE("/:id", controllers.DeleteStudentV1)
	}

	// 为 V1 提供带版本号的别名
	r.Group("/v1/students").Any("/*any", func(c *gin.Context) {
		c.Request.URL.Path = "/students" + c.Param("any")
		r.HandleContext(c)
	})

	// V2: SQLite 原生 SQL
	v2 := r.Group("/v2/students")
	{
		v2.POST("", controllers.CreateStudentV2)
		v2.GET("", controllers.GetStudentsV2)
		// v2.GET("/:id", controllers.GetStudentV2)
		// v2.PUT("/:id", controllers.UpdateStudentV2)
		// v2.DELETE("/:id", controllers.DeleteStudentV2)
	}

	// V3: MySQL + GORM
	v3 := r.Group("/v3/students")
	{
		v3.POST("", controllers.CreateStudentV3)
		v3.GET("", controllers.GetStudentsV3)
		v3.DELETE("/:id", controllers.DeleteStudentV3)
	}
}
