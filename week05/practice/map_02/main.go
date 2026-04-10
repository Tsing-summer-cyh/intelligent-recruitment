package main

import (
	"fmt"
)

func main() {
	// 定义一个测试用的字母数字字符串
	inputStr := "Go2023Golang2024Gooo"

	// 调用函数获取结果
	mostFreqChar, count := findMostFrequentChar(inputStr)

	// 打印结果
	fmt.Printf("输入字符串: %s\n", inputStr)
	if count > 0 {
		fmt.Printf("出现次数最多的字符是: '%c', 一共出现了 %d 次\n", mostFreqChar, count)
	} else {
		fmt.Println("输入的字符串为空！")
	}
}

// findMostFrequentChar 统计字符串中字符的出现次数，并返回出现次数最多的字符及其次数
func findMostFrequentChar(s string) (rune, int) {
	// 边界条件处理：如果是空字符串，直接返回
	if len(s) == 0 {
		return 0, 0
	}

	// 1. 初始化一个 map，键(key)为字符(rune类型)，值(value)为出现的次数(int类型)
	charCount := make(map[rune]int)

	// 用于记录当前出现次数最多的字符和对应的次数
	var maxChar rune
	maxCount := 0

	// 2. 遍历字符串
	for _, char := range s {
		// 将字符的出现次数累加
		charCount[char]++

		// 3. 每次累加后，顺便检查是否超过了当前记录的最大次数
		if charCount[char] > maxCount {
			maxCount = charCount[char] // 更新最大次数
			maxChar = char             // 更新出现次数最多的字符
		}
	}

	return maxChar, maxCount
}