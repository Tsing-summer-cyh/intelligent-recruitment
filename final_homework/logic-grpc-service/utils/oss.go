package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"gopkg.in/yaml.v3"
)

// Config 映射 yaml 配置
type Config struct {
	OSS struct {
		Endpoint        string `yaml:"endpoint"`
		AccessKeyID     string `yaml:"access_key_id"`
		AccessKeySecret string `yaml:"access_key_secret"`
		BucketName      string `yaml:"bucket_name"`
	} `yaml:"oss"`
}

var (
	ossClient *oss.Client
	bucket    *oss.Bucket
	AppConfig Config
)

// InitOSS 初始化 OSS 客户端
func InitOSS() {
	// 1. 读取 YAML 配置文件
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("读取 config.yaml 失败: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &AppConfig)
	if err != nil {
		log.Fatalf("解析 YAML 失败: %v", err)
	}

	// 2. 初始化阿里云 OSS 客户端
	ossClient, err = oss.New(AppConfig.OSS.Endpoint, AppConfig.OSS.AccessKeyID, AppConfig.OSS.AccessKeySecret)
	if err != nil {
		log.Fatalf("创建 OSS 客户端失败: %v", err)
	}

	// 3. 获取存储空间
	bucket, err = ossClient.Bucket(AppConfig.OSS.BucketName)
	if err != nil {
		log.Fatalf("获取 Bucket 失败: %v", err)
	}

	fmt.Println("☁️  私有 OSS 服务连接成功！")
}

// GeneratePresignedURL 生成限时访问的签名 URL (核心合规要求)
// objectName 是简历在 OSS 里的相对路径，比如 "resumes/user_1_resume.pdf"
func GeneratePresignedURL(objectName string) (string, error) {
	if objectName == "" {
		return "", nil
	}
	// 生成有效时间为 10 分钟 (600秒) 的 GET 请求签名 URL
	signedURL, err := bucket.SignURL(objectName, oss.HTTPGet, 600)
	if err != nil {
		return "", fmt.Errorf("生成签名URL失败: %v", err)
	}
	return signedURL, nil
}

// UploadResume 上传简历到 OSS
// objectName 是存储路径，content 是文件内容
func UploadResume(objectName string, content []byte) error {
	if bucket == nil {
		return fmt.Errorf("OSS Bucket 未初始化")
	}
	// 使用 PutObject 上传字节流
	err := bucket.PutObject(objectName, bytes.NewReader(content))
	if err != nil {
		return fmt.Errorf("上传文件失败: %v", err)
	}
	return nil
}