package main

import (
	"math"
	"testing"
)

// TestCalculateFinalFee 测试整体计费逻辑
func TestCalculateFinalFee(t *testing.T) {
	// 利用结构体和数组切片定义测试用例集合
	tests := []struct {
		name        string  // 测试用例名称
		usage       float64 // 测试用电量
		timeStr     string  // 测试时段
		expectedFee float64 // 期望计算结果
	}{
		{
			name:        "第一档_高峰时段",
			usage:       100.0,
			timeStr:     "14:00",
			expectedFee: 55.0, // 100 * 0.5 * 1.1 = 55.0
		},
		{
			name:        "第一档_低谷时段",
			usage:       100.0,
			timeStr:     "02:00",
			expectedFee: 40.0, // 100 * 0.5 * 0.8 = 40.0
		},
		{
			name:        "第二档_高峰时段",
			usage:       300.0,
			timeStr:     "20:00",
			expectedFee: 198.0, // (200*0.5 + 100*0.8) * 1.1 = 180 * 1.1 = 198.0
		},
		{
			name:        "第二档_低谷时段",
			usage:       400.0,
			timeStr:     "23:00",
			expectedFee: 208.0, // (200*0.5 + 200*0.8) * 0.8 = 260 * 0.8 = 208.0
		},
		{
			name:        "第三档_高峰时段",
			usage:       500.0,
			timeStr:     "10:00",
			expectedFee: 418.0, // (200*0.5 + 200*0.8 + 100*1.2) * 1.1 = 380 * 1.1 = 418.0
		},
		{
			name:        "第三档_低谷时段",
			usage:       600.0,
			timeStr:     "05:30",
			expectedFee: 400.0, // (200*0.5 + 200*0.8 + 200*1.2) * 0.8 = 500 * 0.8 = 400.0
		},
	}

	// 遍历执行测试用例
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actualFee := CalculateFinalFee(tc.usage, tc.timeStr)
			
			// 浮点数比较会有精度问题，通常计算两者差值的绝对值是否小于一个极小值
			if math.Abs(actualFee-tc.expectedFee) > 0.001 {
				t.Errorf("用例 [%s] 失败: 输入用电量 %.2f, 时段 %s. 期望电费 %.2f, 实际计算结果 %.2f",
					tc.name, tc.usage, tc.timeStr, tc.expectedFee, actualFee)
			}
		})
	}
}