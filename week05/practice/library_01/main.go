package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {
	// 1. 定义命令行参数
	// flag.String 返回的是一个指针
	filePath := flag.String("file", "", "指定要处理的文件路径 (必填)")
	operation := flag.String("operation", "", "指定操作类型，支持: count, upper, lower (必填)")
	outputPath := flag.String("output", "", "指定输出文件路径 (可选，不填则输出到终端)")

	// 2. 解析命令行参数
	flag.Parse()

	// 3. 参数校验：校验必填参数是否提供
	if *filePath == "" || *operation == "" {
		fmt.Println("❌ 错误：缺少必要的参数。")
		fmt.Println("👉 用法示例：go run main.go -file input.txt -operation count -output result.txt")
		flag.Usage() // 打印默认的帮助信息
		return
	}

	// 4. 读取源文件
	// os.ReadFile 会一次性读取整个文件内容，适用于小文件
	data, err := os.ReadFile(*filePath)
	if err != nil {
		fmt.Printf("❌ 读取文件失败: %v\n", err)
		return
	}
	content := string(data)
	var resultText string

	// 5. 根据 -operation 执行不同的业务逻辑
	switch *operation {
	case "count":
		// 统计字符数（使用 utf8.RuneCountInString 以正确统计中文字符）
		count := utf8.RuneCountInString(content)
		resultText = fmt.Sprintf("文件 '%s' 共包含 %d 个字符。\n", *filePath, count)
	case "upper":
		// 将文本转为大写
		resultText = strings.ToUpper(content)
	case "lower":
		// 将文本转为小写
		resultText = strings.ToLower(content)
	default:
		fmt.Printf("❌ 不支持的操作类型: '%s'。目前仅支持: count, upper, lower\n", *operation)
		return
	}

	// 6. 处理输出
	if *outputPath != "" {
		// 如果指定了 -output 参数，则将结果写入目标文件
		// 0644 表示文件权限：所有者可读写，其他人可读
		err := os.WriteFile(*outputPath, []byte(resultText), 0644)
		if err != nil {
			fmt.Printf("❌ 写入输出文件失败: %v\n", err)
			return
		}
		fmt.Printf("✅ 操作成功！结果已保存至: %s\n", *outputPath)
	} else {
		// 如果没有指定 -output，则直接打印到控制台
		fmt.Println("====== 执行结果 ======")
		fmt.Print(resultText)
		fmt.Println("\n======================")
	}
}