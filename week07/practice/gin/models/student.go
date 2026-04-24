package models

// Student 学生信息结构体
type Student struct {
	ID    int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name  string `json:"name" binding:"required"`
	Age   int    `json:"age" binding:"required,gt=0"`
	Grade string `json:"grade" binding:"required"`
}
