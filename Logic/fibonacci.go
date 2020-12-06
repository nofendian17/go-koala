package main

import "fmt"

// @todo : Fibonacci with range 1-1000

func fibonacci(n int) []int {
	var a = 0
	var b = 0
	var sum = 1
	var res []int

	for b <= n {
		a = b
		b = sum
		sum = a + b
		res = append(res, a)
	}
	return res
}

func main() {
	input := 1000
	result := fibonacci(input)
	fmt.Println(result)
}
