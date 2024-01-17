package main

import (
	"fmt"
	"math"
)

func main() {
	nums := []int{1, 3}
	k := 3
	fmt.Println(countKDifference(nums, k))
}

func countKDifference(nums []int, k int) int {
	var resp int
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			fmt.Println(i, j)
			if math.Abs(float64(nums[i])-float64(nums[j])) == float64(k) {
				resp++
			}
		}
	}

	fmt.Println(resp)

	return 0
}
