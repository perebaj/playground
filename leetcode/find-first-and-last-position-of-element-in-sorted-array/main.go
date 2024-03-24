package main

import (
	"sort"
)

func main() {
	searchRange([]int{2, 2}, 3)
}

func searchRange(nums []int, target int) []int {
	if len(nums) == 0 {
		return []int{-1, -1}
	}

	upperBound := sort.Search(len(nums), func(i int) bool {
		return nums[i] > target
	})

	lowerBound := sort.Search(len(nums), func(i int) bool {
		return nums[i] >= target
	})
	if lowerBound == len(nums) || nums[lowerBound] != target {
		return []int{-1, -1}
	}

	return []int{lowerBound, upperBound - 1}
}
