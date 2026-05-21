package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/joho/godotenv"
)

var (
	ossClient *oss.Client
	bucket    *oss.Bucket
)

// InitOSS 初始化 OSS 客户端（从 .env 读取配置）
func InitOSS() {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ 未找到 .env 文件，将使用系统环境变量")
	}

	// 从环境变量读取配置
	endpoint := os.Getenv("OSS_ENDPOINT")
	accessKeyID := os.Getenv("OSS_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("OSS_ACCESS_KEY_SECRET")
	bucketName := os.Getenv("OSS_BUCKET_NAME")

	// 校验必要配置
	if endpoint == "" || accessKeyID == "" || accessKeySecret == "" || bucketName == "" {
		log.Fatalf("❌ OSS 配置缺失，请检查 .env 文件中的 OSS_* 环境变量")
	}

	// 初始化阿里云 OSS 客户端
	ossClient, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		log.Fatalf("创建 OSS 客户端失败: %v", err)
	}

	// 获取存储空间
	bucket, err = ossClient.Bucket(bucketName)
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