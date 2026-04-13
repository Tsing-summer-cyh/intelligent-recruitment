package main

import (
	"fmt"
	"os"
	"time"
)

// GenerateReport 聚合数据并写入日志
func GenerateReport(results []ProbeResult) {
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("monitor-log-%s.log", timestamp)

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("无法创建日志文件: %v\n", err)
		return
	}
	defer file.Close()

	total := len(results)
	success := 0
	var maxLatency time.Duration
	var slowestTarget Target

	// 统计核心指标
	for _, r := range results {
		if r.Status {
			success++
		}
		// 寻找耗时最长的请求
		if r.ResponseTime > maxLatency {
			maxLatency = r.ResponseTime
			slowestTarget = r.Target
		}
	}

	fail := total - success
	var successRate float64
	if total > 0 {
		successRate = float64(success) / float64(total) * 100
	}

	// 构建报表文本
	report := "=== 服务健康探测分析报告 ===\n"
	report += fmt.Sprintf("生成时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	report += fmt.Sprintf("探测总数: %d\n", total)
	report += fmt.Sprintf("成功数量: %d\n", success)
	report += fmt.Sprintf("失败数量: %d\n", fail)
	report += fmt.Sprintf("成功比率: %.2f%%\n", successRate)
	if total > 0 {
		report += fmt.Sprintf("最慢服务: %s (耗时: %v)\n", slowestTarget.Name, maxLatency)
	}
	report += "--------------------------------------\n"
	report += "详细结果:\n"

	for _, r := range results {
		statusStr := "【成功】"
		if !r.Status {
			statusStr = fmt.Sprintf("【失败】(%s)", r.ErrorMsg)
		}
		report += fmt.Sprintf("- %s | 目标: %s | 耗时: %v | 状态: %s\n", r.Target.Name, r.Target.URL, r.ResponseTime, statusStr)
	}

	// 写入文件
	file.WriteString(report)
	fmt.Printf("\n[系统提示] 探测完成！报表已生成: %s\n", filename)
}