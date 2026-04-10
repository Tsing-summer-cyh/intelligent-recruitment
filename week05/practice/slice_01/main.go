package main

import (
	"fmt"
)

func main() {
	// 1. 创建一个包含 10 个整数的切片，初始值为 1 到 10
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Println("1. 初始切片:", slice)

	// 2. 使用切片操作获取第 3 到第 7 个元素（包含第 7 个）
	// 说明：Go 语言索引从 0 开始。第 3 个元素索引为 2，第 7 个元素索引为 6。
	// 切片操作 [start:end] 是左闭右开区间，因此截取 [2:7] 刚好包含索引 2, 3, 4, 5, 6。
	subSlice := slice[2:7]
	fmt.Println("2. 第 3 到第 7 个元素:", subSlice)

	// 3. 在切片末尾添加三个新元素：11, 12, 13
	// 说明：使用内置函数 append，这可能会触发切片的底层数组扩容。
	slice = append(slice, 11, 12, 13)
	fmt.Println("3. 添加 11, 12, 13 后的切片:", slice)

	// 4. 删除切片中的第 5 个元素
	// 说明：第 5 个元素的索引为 4。在 Go 中，通常通过 append 拼接该元素前后的切片段来实现删除。
	// slice[:4] 取出前 4 个元素，slice[5:] 取出第 5 个元素之后的所有元素，然后将它们解包 (...) 追加在一起。
	slice = append(slice[:4], slice[5:]...)
	fmt.Println("4. 删除第 5 个元素后的切片:", slice)

	// 5. 将切片中的所有元素乘以 2
	// 说明：使用 for range 遍历切片的索引，直接修改原切片的值。
	for i := range slice {
		slice[i] *= 2
	}
	fmt.Println("5. 所有元素乘以 2 后的切片:", slice)

	// 6. 打印最终切片的内容和容量
	// 说明：使用 cap() 函数获取切片的容量。
	fmt.Printf("6. 最终切片内容: %v\n", slice)
	fmt.Printf("   最终切片容量: %d\n", cap(slice))
}