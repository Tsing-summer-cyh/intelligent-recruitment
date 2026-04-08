package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// 使用常量 (const) 定义所有的阶梯阈值和单价
const (
	Tier1Limit = 200.0 // 第一档用电量上限
	Tier2Limit = 400.0 // 第二档用电量上限

	Tier1Price = 0.5 // 第一档单价 (0-200度)
	Tier2Price = 0.8 // 第二档单价 (200-400度超出部分)
	Tier3Price = 1.2 // 第三档单价 (400度以上超出部分)

	PeakRate   = 1.1 // 高峰时段费率 (+10%)
	ValleyRate = 0.8 // 低谷时段费率 (-20%)
)

// 使用 init() 函数初始化系统
func init() {
	fmt.Println("------------------------------------------------")
	fmt.Println("计费规则版本号: v1.0.0")
	fmt.Printf("系统初始化时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("------------------------------------------------")
}

// CalculateBaseFee 计算基础阶梯电费
func CalculateBaseFee(usage float64) float64 {
	if usage <= 0 {
		return 0.0
	}

	if usage <= Tier1Limit {
		return usage * Tier1Price
	} else if usage <= Tier2Limit {
		// 满第一档的费用 + 超出第一档部分的费用
		return (Tier1Limit * Tier1Price) + ((usage - Tier1Limit) * Tier2Price)
	} else {
		// 满第一档的费用 + 满第二档的费用 + 超出第二档部分的费用
		return (Tier1Limit * Tier1Price) + ((Tier2Limit - Tier1Limit) * Tier2Price) + ((usage - Tier2Limit) * Tier3Price)
	}
}

// IsPeakTime 判断是否为高峰时段
// 假设输入格式为 "HH:mm"，高峰定义为 8:00 到 21:59 (即 8 <= hour < 22)
func IsPeakTime(timeStr string) bool {
	parts := strings.Split(timeStr, ":")
	if len(parts) > 0 {
		hour, err := strconv.Atoi(parts[0])
		if err == nil {
			if hour >= 8 && hour < 22 {
				return true
			}
		}
	}
	return false // 解析失败或不在区间内，均算作低谷
}

// CalculateFinalFee 根据时段计算最终电费
func CalculateFinalFee(usage float64, timeStr string) float64 {
	baseFee := CalculateBaseFee(usage)

	if IsPeakTime(timeStr) {
		return baseFee * PeakRate
	}
	return baseFee * ValleyRate
}

// RunBillingSystem 执行交互式计费流程
func RunBillingSystem() {
	var usage float64
	var timeStr string

	// 引导用户输入
	fmt.Print("请输入用户用电量(比如 400.00): ")
	fmt.Scanln(&usage)

	fmt.Print("请输入用户用电时段(比如 14:00): ")
	fmt.Scanln(&timeStr)

	// 计算费用
	finalFee := CalculateFinalFee(usage, timeStr)

	// 打印账单
	fmt.Println("\n--- 账单明细 ---")
	fmt.Printf("当前用电：%.2f 度\n", usage)
	fmt.Printf("当前时段：%s 点\n", timeStr)
	fmt.Printf("最终电费：%.2f 元\n", finalFee)
}

func main() {
	RunBillingSystem()
}