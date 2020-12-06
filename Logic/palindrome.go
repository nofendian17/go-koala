package main

import (
	"fmt"
	"strings"
)

// todo : check palindrome string
// ex : maddam result is true

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func isPalindrome(word string) bool {
	if word == reverse(word) {
		return true
	}
	return false
}

func main() {
	var str string
	fmt.Print("Enter a string: ")
	fmt.Scan(&str)
	if isPalindrome(strings.ToUpper(str)) == true {
		fmt.Print(str, " is a palindrome.")
	} else {
		fmt.Print(str, " is not a palindrome.")
	}
}
