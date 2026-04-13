package main

import (
	"encoding/json"
	"os"
	"time"
)

// Target 定义单个探测目标
type Target struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Expect     string `json:"expect"`
	RetryCount int    `json:"retry_count"`
}

// Config 对应 JSON 的根节点
type Config struct {
	Targets []Target `json:"targets"`
}

// ProbeResult 存储单个目标的最终探测结果
type ProbeResult struct {
	Target       Target
	Status       bool
	ResponseTime time.Duration
	ErrorMsg     string
}

// LoadConfig 从指定路径加载并解析配置文件
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}