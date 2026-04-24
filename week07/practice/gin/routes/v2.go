package routes

import (
	v2 "week07/practice/gin/controllers/v2"

	"github.com/gin-gonic/gin"
)

func RegisterV2Routes(r *gin.Engine) {
	v2Group := r.Group("/v2/students")
	{
		v2Group.POST("", v2.CreateStudent)
		v2Group.GET("", v2.GetStudents)
		v2Group.GET("/:id", v2.GetStudent)
		v2Group.PUT("/:id", v2.UpdateStudent)
		v2Group.DELETE("/:id", v2.DeleteStudent)
	}
}
