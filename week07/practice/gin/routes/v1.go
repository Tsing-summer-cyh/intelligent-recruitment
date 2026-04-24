package routes

import (
	v1 "week07/practice/gin/controllers/v1"

	"github.com/gin-gonic/gin"
)

func RegisterV1Routes(r *gin.Engine) {
	v1Group := r.Group("/v1/students")
	{
		v1Group.POST("", v1.CreateStudent)
		v1Group.GET("", v1.GetStudents)
		v1Group.GET("/:id", v1.GetStudent)
		v1Group.PUT("/:id", v1.UpdateStudent)
		v1Group.DELETE("/:id", v1.DeleteStudent)
	}
}
