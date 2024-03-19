package main

import "fmt"

func main() {
	nums := []int{-7, 2, -3, 8, 9, 0, 8, 4, -5, 8, -5, -5, 1, 6, 4}
	fmt.Println(longestConsecutive(nums))
}

func longestConsecutive(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	m := make(map[int]int)
	for _, v := range nums {
		m[v]++
	}
	var initalElementSlice []int
	for k := range m {
		initial := k - 1
		_, ok := m[initial]
		if !ok {
			initalElementSlice = append(initalElementSlice, k)
		}
	}
	positionInitialElementSlice := 0
	var result int
	curResult := 1
	for {
		_, ok := m[initalElementSlice[positionInitialElementSlice]+curResult]
		if !ok {
			if curResult > result {
				result = curResult
			}
			curResult = 1
			positionInitialElementSlice++
		} else {
			curResult++
		}

		if positionInitialElementSlice >= len(initalElementSlice) {
			break
		}
	}

	return result
}
