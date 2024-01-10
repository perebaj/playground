package main

import "fmt"

func main() {
	nums := []int{1, 2, 1}
	fmt.Println(getConcatenation(nums))
}

func getConcatenation(nums []int) []int {
	var resp []int
	for i := 0; i < 2; i++ {
		for _, v := range nums {
			resp = append(resp, v)
		}
	}

	return resp
}
