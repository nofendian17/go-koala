package main

import "fmt"

// @todo : print prime number range 1-1000

func prime(n int) []int {
	var res []int
	for i := 1; i <= n; i++ {
		var counter = 0
		for j := 1; j <= i; j++ {
			if i%j == 0 {
				counter++
			}
		}

		if counter == 2 {
			res = append(res, i)
		}
	}
	return res
}

func main() {
	input := 1000
	result := prime(input)
	fmt.Println(result)
}
