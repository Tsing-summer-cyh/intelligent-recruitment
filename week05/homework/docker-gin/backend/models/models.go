package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"unique;not null" json:"username"`
	PasswordHash string         `gorm:"not null" json:"-"` // 不返回给前端
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// SentenceList 用于处理 MySQL 中的 JSON 数组
type SentenceList []string

func (s *SentenceList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, s)
}

func (s SentenceList) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Word 单词模型
type Word struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"index" json:"user_id"`
	Word       string         `gorm:"not null" json:"word"`
	Meaning    string         `gorm:"type:text" json:"meaning"`
	Sentences  SentenceList   `gorm:"type:json" json:"sentences"`
	AIProvider string         `gorm:"column:ai_provider" json:"ai_provider"` // 👈 只需在这里加上 column:ai_provider
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
