package main

import (
	"fmt"
)

func main() {
	// 1. 定义两个初始切片
	slice1 := []int{1, 2, 3, 4}
	slice2 := []int{3, 4, 5, 6}

	// 2. 将这两个切片拼接成一个新的切片 combinedSlice
	// 说明：使用 append 函数，并通过 `...` 将 slice2 解包成独立的元素追加到 slice1 后面
	combinedSlice := append(slice1, slice2...)
	fmt.Println("拼接后的 combinedSlice:", combinedSlice)

	// 3. 对 combinedSlice 进行去重操作，得到 uniqueSlice
	// 说明：利用 map 的键不能重复的特性来进行去重
	seen := make(map[int]bool) // 用于记录某个元素是否已经出现过
	var uniqueSlice []int      // 用于存放去重后的新切片

	for _, num := range combinedSlice {
		// 如果 seen 这个 map 中不存在 num 这个键（即值为 false）
		if !seen[num] {
			seen[num] = true // 标记为已出现
			uniqueSlice = append(uniqueSlice, num) // 将不重复的元素加入新切片
		}
	}

	// 4. 打印最终的 uniqueSlice
	fmt.Println("去重后的 uniqueSlice:", uniqueSlice)
}