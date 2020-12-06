package helpers

import (
	"strconv"
	"time"
)

func GenerateLetter(n int) string {
	var preffix string = "PO"
	var number string = strconv.Itoa(n)
	var romawian string = RomawianChar(CurrentMonth())
	var currentYear = strconv.Itoa(CurrentYear())

	return preffix + "-" + number + "/" + romawian + "/" + currentYear
}

func RomawianChar(n int) string {
	decimal := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	romawi := []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}

	var result string

	for i := 0; i < len(decimal); i++ {
		for decimal[i] <= n {
			result += romawi[i]
			n -= decimal[i]
		}
	}
	return result
}

func CurrentYear() int {
	t := time.Now()
	year := t.Year()
	return year
}

func CurrentMonth() int {
	t := time.Now()
	month := t.Month()
	return int(month)
}
