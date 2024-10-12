package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	// 已处理输入参数
	args := os.Args[1:]
	a, _ := strconv.Atoi(args[0])
	b, _ := strconv.Atoi(args[1])
	fmt.Println(SumOfTwoNumbers(a, b))
}

func SumOfTwoNumbers(a, b int) int {
	// 存在死循环
	for {
		a++
	}
	// 解题代码请写于此处：
	return a + b
}
