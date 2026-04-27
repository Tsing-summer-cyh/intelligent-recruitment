package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type AIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AIResponse 定义大模型外层返回结构
type AIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// FetchAIResult 调用 AI 接口并解析结果
func FetchAIResult(word string, provider string) (string, []string, error) {
	// ⚠️ 注意：这里必须填入你真实申请的 DeepSeek API Key
	apiKey := "sk-ac246b17607143df915bf97e05243aba"
	url := "https://api.deepseek.com/chat/completions"

	// 优化 Prompt，强制 AI 只输出纯粹的 JSON
	prompt := fmt.Sprintf("请详细解释单词 '%s'。你必须严格返回合法的 JSON 格式，不要包含任何额外的说明文字或 Markdown 标记。格式如下: {\"meaning\": \"中文释义\", \"sentences\": [\"例句1\", \"例句2\", \"例句3\"]}", word)

	reqBody, _ := json.Marshal(AIRequest{
		Model: "deepseek-chat", // 使用 deepseek 的模型名称
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
	})

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	// 1. 解析外层 HTTP 响应数据
	var aiResp AIResponse
	if err := json.Unmarshal(bodyBytes, &aiResp); err != nil || len(aiResp.Choices) == 0 {
		return "", nil, fmt.Errorf("请求大模型失败，可能 Key 无效或网络异常")
	}

	// 提取 AI 实际生成的文本内容
	content := aiResp.Choices[0].Message.Content

	// 2. 清理 AI 可能带上的 markdown 代码块标记 (非常关键的容错处理)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	// 3. 解析内层业务需要的 JSON 数据
	var result struct {
		Meaning   string   `json:"meaning"`
		Sentences []string `json:"sentences"`
	}

	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return "", nil, fmt.Errorf("AI 返回格式异常无法解析")
	}

	return result.Meaning, result.Sentences, nil
}