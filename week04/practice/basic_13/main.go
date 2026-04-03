package main

import (
	"fmt"
)

// Divide 接受两个浮点数并执行除法。
// 如果除数为零，则使用 panic 抛出异常。
func Divide(a, b float64) float64 {
	if b == 0 {
		panic("除数不能为零")
	}
	return a / b
}

// CallDivide 是调用 Divide 的包装函数。
// 它通过 recover 捕获可能发生的 panic，并返回一个友好的错误信息。
func CallDivide(a, b float64) (result float64, err error) {
	// 必须在 defer 函数中调用 recover
	defer func() {
		if r := recover(); r != nil {
			// 捕获到 panic，将其格式化为 error 并赋值给命名返回值 err
			err = fmt.Errorf("捕获到除法异常: %v", r)
		}
	}()

	// 执行除法，如果这里发生 panic，程序会直接跳转到 defer 中执行
	result = Divide(a, b)
	
	// 如果没有异常，正常返回
	return result, nil
}

func main() {
	// 测试用例 1：正常的除法操作
	fmt.Println("--- 测试正常情况 ---")
	res1, err1 := CallDivide(10.0, 2.0)
	if err1 != nil {
		fmt.Println("错误:", err1)
	} else {
		fmt.Printf("10.0 / 2.0 = %.2f\n", res1)
	}

	fmt.Println()

	// 测试用例 2：触发 panic 的除法操作
	fmt.Println("--- 测试除数为零情况 ---")
	res2, err2 := CallDivide(5.0, 0.0)
	if err2 != nil {
		fmt.Println("错误:", err2)
	} else {
		fmt.Printf("5.0 / 0.0 = %.2f\n", res2)
	}
}