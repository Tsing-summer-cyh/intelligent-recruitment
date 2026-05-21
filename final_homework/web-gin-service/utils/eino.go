package utils

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

var chatModel model.ChatModel

// InitEino 初始化大模型客户端
func InitEino() {
	// ⚠️ 将这里替换为你刚刚在阿里云百炼申请的真实 API Key
	apiKey := "sk-b7b51dc1e6044449871be581c07c78b3"

	config := &openai.ChatModelConfig{
		// 百炼提供的 OpenAI 兼容访问地址
		BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
		APIKey:  apiKey,
		Model:   "qwen-plus", // 选用通义千问 Plus 模型
	}

	var err error
	chatModel, err = openai.NewChatModel(context.Background(), config)
	if err != nil {
		fmt.Printf("❌ Eino 初始化失败: %v\n", err)
		return
	}
	fmt.Println("🧠 Eino AI 引擎初始化成功，已连接通义千问大模型！")
}

// ChatWithAI 发送对话请求
func ChatWithAI(ctx context.Context, question string) (string, error) {
	if chatModel == nil {
		return "", fmt.Errorf("AI 模型未初始化")
	}

	// 构造 Eino 的消息体
	msg := schema.UserMessage(question)

	// 调用大模型生成回复
	resp, err := chatModel.Generate(ctx, []*schema.Message{msg})
	if err != nil {
		return "", err
	}

	return resp.Content, nil
}