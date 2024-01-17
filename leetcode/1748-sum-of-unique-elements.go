package main

import "fmt"

func main() {
	nums := []int{1, 2, 3, 2}
	fmt.Println(sumOfUnique(nums))
}

func sumOfUnique(nums []int) int {
	m := make(map[int]int)

	for _, v := range nums {
		m[v]++
	}

	var response int
	for k, v := range m {
		if v == 1 {
			response = response + k
		}
	}

	return response
}
