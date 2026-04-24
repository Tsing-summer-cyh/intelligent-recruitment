package routes

import (
	v4 "week07/practice/gin/controllers/v4"

	"github.com/gin-gonic/gin"
)

func RegisterV4Routes(r *gin.Engine) {
	v4Group := r.Group("/v4/students")
	{
		v4Group.POST("", v4.CreateStudent)
		v4Group.GET("", v4.GetStudents)
		v4Group.GET("/:id", v4.GetStudent)
		v4Group.PUT("/:id", v4.UpdateStudent)
		v4Group.DELETE("/:id", v4.DeleteStudent)
	}
}
