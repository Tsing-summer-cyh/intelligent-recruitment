package models

// Student 学生信息结构体
type Student struct {
	ID    int    `json:"id" gorm:"primaryKey;autoIncrement"` // 学生 ID
	Name  string `json:"name" binding:"required"`            // 学生姓名
	Age   int    `json:"age" binding:"required,gt=0"`        // 学生年龄
	Grade string `json:"grade" binding:"required"`           // 学生年级
}
