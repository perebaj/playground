package main

import (
	"fmt"
	"sort"
)

func main() {
	nums := []int{3, 2, 2, 3}
	removeElement(nums, 3)
}

func removeElement(nums []int, val int) int {
	var kResult int
	for i, v := range nums {
		if v == val {
			kResult++
			nums[i] = 101
		}
	}
	kResult = len(nums) - kResult
	sort.Ints(nums)
	fmt.Println(nums)
	for i, v := range nums {
		if v == 101 {
			nums[i] = 0
		}
	}
	fmt.Println(nums)

	return kResult
}
