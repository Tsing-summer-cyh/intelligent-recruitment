package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

var chatModel model.ChatModel

// InitEino 初始化大模型客户端（从 .env 读取配置）
func InitEino() {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ 未找到 .env 文件，将使用系统环境变量")
	}

	// 从环境变量读取配置
	apiKey := os.Getenv("DASHSCOPE_API_KEY")
	modelName := os.Getenv("DASHSCOPE_MODEL")

	// 校验必要配置
	if apiKey == "" {
		log.Fatalf("❌ AI 配置缺失，请检查 .env 文件中的 DASHSCOPE_API_KEY")
	}
	if modelName == "" {
		modelName = "qwen-plus" // 默认模型
	}

	config := &openai.ChatModelConfig{
		BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
		APIKey:  apiKey,
		Model:   modelName,
	}

	var err error
	chatModel, err = openai.NewChatModel(context.Background(), config)
	if err != nil {
		fmt.Printf("❌ Eino 初始化失败: %v\n", err)
		return
	}
	fmt.Println("🧠 Eino AI 引擎初始化成功，已连接通义千问大模型！")
}

// ChatWithAI 发送对话请求（带数据库上下文）
func ChatWithAI(ctx context.Context, question string, dbContext string) (string, error) {
	if chatModel == nil {
		return "", fmt.Errorf("AI 模型未初始化")
	}

	// 构建系统提示词，告知 AI 它是一个招聘系统助手，有真实数据
	systemPrompt := `你是一个智能招聘系统的 AI 助手，专门帮助 HR 管理招聘工作。
你的职责是：
1. 回答关于招聘岗位、投递情况、候选人信息的查询
2. 帮助 HR 分析招聘数据，提供见解和建议
3. 起草招聘相关的文档（如岗位描述、面试问题等）

当回答涉及具体数据的问题时，请优先使用系统提供的真实数据库数据来回答，不要编造数据。
如果用户问的数据你没有相关信息，请诚实告知用户当前没有相关数据。

回答风格：专业、简洁、有帮助，避免过长的废话。`

	// 构建消息列表
	messages := []*schema.Message{
		schema.SystemMessage(systemPrompt),
	}

	// 如果有数据库上下文，添加数据信息
	if dbContext != "" {
		dataPrompt := fmt.Sprintf("以下是系统数据库中的真实数据信息，请在回答时参考这些数据：\n\n%s", dbContext)
		messages = append(messages, schema.UserMessage(dataPrompt))
		// 添加助手确认消息（使用字符串 "assistant" 作为角色）
		messages = append(messages, &schema.Message{
			Role:    "assistant",
			Content: "好的，我已经了解了当前的数据库数据情况，会基于这些真实数据来回答您的问题。",
		})
	}

	// 添加用户的实际问题
	messages = append(messages, schema.UserMessage(question))

	resp, err := chatModel.Generate(ctx, messages)
	if err != nil {
		return "", err
	}

	return resp.Content, nil
}

// DetectIntent 检测用户问题的意图类型
func DetectIntent(question string) string {
	question = strings.ToLower(question)

	// 岗位相关
	if strings.Contains(question, "岗位") || strings.Contains(question, "职位") ||
		strings.Contains(question, "招聘") || strings.Contains(question, "发布") {
		return "jobs"
	}

	// 投递相关
	if strings.Contains(question, "投递") || strings.Contains(question, "申请") ||
		strings.Contains(question, "简历") || strings.Contains(question, "应聘") {
		return "applications"
	}

	// 候选人相关
	if strings.Contains(question, "候选人") || strings.Contains(question, "求职者") ||
		strings.Contains(question, "应聘者") || strings.Contains(question, "档案") {
		return "candidates"
	}

	// 统计相关
	if strings.Contains(question, "统计") || strings.Contains(question, "数量") ||
		strings.Contains(question, "多少") || strings.Contains(question, "数据") ||
		strings.Contains(question, "分析") || strings.Contains(question, "趋势") {
		return "stats"
	}

	// 默认
	return "general"
}