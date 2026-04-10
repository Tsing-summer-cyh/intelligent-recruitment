package main

import (
	"fmt"
	"sync"
	"time"
)

// calcSquareSum 负责计算切片片段中所有元素的平方和
func calcSquareSum(chunk []int, resultChan chan<- int, wg *sync.WaitGroup) {
	// 确保在函数退出时通知 WaitGroup 计数器减 1
	defer wg.Done()

	sum := 0
	for _, num := range chunk {
		sum += num * num
		// 题目要求的 time.Sleep(1)
		// 注意：在 Go 语言中 time.Sleep 接收的是 time.Duration 类型，
		// 直接传 1 代表 1 纳秒 (Nanosecond)。这里用来模拟计算过程中的微小耗时。
		time.Sleep(1)
	}

	// 将计算出的局部平方和发送到通道中
	resultChan <- sum
}

func main() {
	// 1. 创建并初始化一个长度为 1000 的整数切片
	totalElements := 1000
	nums := make([]int, totalElements)
	for i := 0; i < totalElements; i++ {
		nums[i] = i + 1 // 初始化为 1 到 1000
	}

	// 2. 定义并发策略
	numGoroutines := 10 // 将切片分为 10 个部分，启动 10 个 Goroutine
	chunkSize := totalElements / numGoroutines

	// 创建通道用于收集各个 Goroutine 的计算结果，带有缓冲可以防止阻塞
	resultChan := make(chan int, numGoroutines)
	
	// 使用 WaitGroup 来等待所有 Goroutine 执行完毕
	var wg sync.WaitGroup

	fmt.Printf("开始启动 %d 个 Goroutine 进行并发计算...\n", numGoroutines)

	// 3. 分割切片并启动 Goroutine
	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := start + chunkSize
		// 处理最后一部分可能不能整除的情况（虽然 1000/10 能整除，但这是严谨的写法）
		if i == numGoroutines-1 {
			end = totalElements
		}

		// 获取当前片段
		chunk := nums[start:end]

		// 增加 WaitGroup 计数
		wg.Add(1)
		
		// 启动 Goroutine
		go calcSquareSum(chunk, resultChan, &wg)
	}

	// 4. 在后台启动一个 Goroutine 等待所有计算完成，然后关闭通道
	go func() {
		wg.Wait()
		close(resultChan) // 所有结果都发送完毕后，关闭通道
	}()

	// 5. 汇总所有 Goroutine 发送过来的结果
	totalSquareSum := 0
	// 遍历通道，直到通道被 close 且数据被读完
	for partialSum := range resultChan {
		totalSquareSum += partialSum
	}

	// 6. 输出最终结果
	fmt.Printf("🎉 所有 Goroutine 计算完毕！\n")
	fmt.Printf("切片元素的最终平方和为: %d\n", totalSquareSum)
}