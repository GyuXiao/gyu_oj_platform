package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// 已处理输入参数
	args := os.Args[1:]
	inputArgs := strings.Split(args[0], " ")
	a, _ := strconv.Atoi(inputArgs[0])
	b, _ := strconv.Atoi(inputArgs[1])
	fmt.Println(SumOfTwoNumbers(a, b))
}

func SumOfTwoNumbers(a, b int) int {
	// 解题代码请写于此处：
	return a + b
}
