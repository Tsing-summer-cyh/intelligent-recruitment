package routes

import (
	v3 "week07/practice/gin/controllers/v3"

	"github.com/gin-gonic/gin"
)

func RegisterV3Routes(r *gin.Engine) {
	v3Group := r.Group("/v3/students")
	{
		v3Group.POST("", v3.CreateStudent)
		v3Group.GET("", v3.GetStudents)
		v3Group.GET("/:id", v3.GetStudent)
		v3Group.PUT("/:id", v3.UpdateStudent)
		v3Group.DELETE("/:id", v3.DeleteStudent)
	}
}
