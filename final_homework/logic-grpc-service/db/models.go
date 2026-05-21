package db

import "time"

// 1. 账号表
type User struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"`
	Username     string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	PasswordHash string    `gorm:"type:varchar(255);not null"` // 实际项目中应存哈希值，这里为演示暂存明文或简单密码
	Role         string    `gorm:"type:varchar(20);not null"`  // "hr" 或 "candidate"
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

// 2. 招聘岗位表
type Job struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	HrID        int64     `gorm:"not null"` // 关联 User.ID
	Title       string    `gorm:"type:varchar(100);not null"`
	Description string    `gorm:"type:text;not null"`
	Status      int32     `gorm:"type:tinyint;default:1"` // 1-上架, 0-下架
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// 3. 候选人结构化档案表
type CandidateProfile struct {
	ID               int64     `gorm:"primaryKey;autoIncrement"`
	UserID           int64     `gorm:"uniqueIndex;not null"` // 关联 User.ID
	RealName         string    `gorm:"type:varchar(50);not null"`
	Phone            string    `gorm:"type:varchar(20);not null"`
	HighestEducation string    `gorm:"type:varchar(50);not null"`
	University       string    `gorm:"type:varchar(100);not null"`
	Experience       string    `gorm:"type:text;not null"`
	Skills           string    `gorm:"type:varchar(255);not null"`
	ResumeOssUrl     string    `gorm:"type:varchar(500)"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`
}

// 4. 岗位投递关联表
type JobApplication struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	JobID       int64     `gorm:"not null"` // 关联 Job.ID
	CandidateID int64     `gorm:"not null"` // 关联 User.ID
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

// 5. AI 对话历史记录表
type AiChatHistory struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	HrID      int64     `gorm:"not null"`
	Question  string    `gorm:"type:text;not null"`
	Answer    string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}