package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// UserRequest 定义了客户端请求的结构体，匹配作业中的请求体格式
type UserRequest struct {
	Stream  int    `json:"stream"`                     // 可选参数，1代表流式输出，不传默认0
	Message string `json:"message" binding:"required"` // 提出的问题，为文本字符串
}

func main() {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("警告: 未找到 .env 文件，将尝试读取系统环境变量")
	}

	// 读取鉴权信息
	authKey := os.Getenv("AUTHORIZATION")
	if authKey == "" {
		log.Fatal("致命错误: 环境变量中未配置 AUTHORIZATION")
	}

	r := gin.Default()

	// 注册 POST 路由 /api/ai/call
	r.POST("/api/ai/call", func(c *gin.Context) {
		var req UserRequest

		// 处理参数解析过程中的错误
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "参数解析失败",
				"details": err.Error(),
			})
			return
		}

		// 根据客户端传入参数判断是否为流式输出 (对应 curl 中的 "stream": true)
		isStream := false
		if req.Stream == 1 {
			isStream = true
		}

		// 构造向秘塔 API 发送的请求体，严格对应 curl 中的 JSON 结构
		metasoReqBody := map[string]interface{}{
			"model": "fast",
			"messages": []map[string]string{
				{"role": "user", "content": req.Message},
			},
			"stream": isStream,
		}

		jsonData, err := json.Marshal(metasoReqBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "请求体打包失败"})
			return
		}

		// 【关键修改点 1】：使用 curl 中提供的最新 API 地址
		apiURL := "https://metaso.cn/api/v1/chat/completions"
		httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建API请求失败"})
			return
		}

		// 【关键修改点 2】：严格按照 curl 补充请求头
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Accept", "application/json") // 新增 Accept 头
		httpReq.Header.Set("Authorization", authKey)

		// 发起 API 调用 (保留了之前绕过系统代理的设置，以防网络环境拦截)
		client := &http.Client{
			Transport: &http.Transport{
				Proxy: nil,
			},
		}
		resp, err := client.Do(httpReq)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error":   "请求秘塔API异常",
				"details": err.Error(),
			})
			return
		}
		defer resp.Body.Close()

		// 响应处理：将 API 响应原样返回给客户端
		for k, values := range resp.Header {
			for _, v := range values {
				c.Writer.Header().Add(k, v)
			}
		}
		c.Writer.WriteHeader(resp.StatusCode)
		io.Copy(c.Writer, resp.Body)
	})

	fmt.Println("服务端启动成功，运行在 http://localhost:8080")
	r.Run(":8080")
}
