package main

import "fmt"

func main() {
	var a, b int
	fmt.Scanln(&a, &b)
	fmt.Println(SumOfTwoNumbers(a, b))
}

func SumOfTwoNumbers(a, b int) int {
	// 解题代码请写于此处：
	return a + b
}
