package main

import "fmt"

func SortedSquaredArray(array []int) []int {
	// -2, 1
	// 2ˆ2 > 1 ^2
	start := 0
	end := len(array) - 1
	result := make([]int, len(array))
	for i := len(array) - 1; i >= 0; i-- {
		if abs(array[start]) > abs(array[end]) {
			result[i] = array[start] * array[start]
			start++
		} else {
			result[i] = array[end] * array[end]
			end--
		}
	}
	fmt.Println(result)
	return result
}

func abs(num int) int {
	if num < 0 {
		return num * -1
	} else {
		return num
	}
}
