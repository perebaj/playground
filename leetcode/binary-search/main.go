package main

import "fmt"

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 8, 9, 9, 11}
	fmt.Println(binarySearch(nums, 2))
}

func binarySearch(nums []int, target int) int {
	left := 0
	right := len(nums) - 1
	index := -1
	for left <= right {
		cur := (right + left) / 2
		if nums[cur] > target {
			//iterage over the left side of my array
			right = cur - 1
		} else if nums[cur] < target {
			left = cur + 1
			//iterate over the right side of my array
		} else {
			return cur
		}
	}

	return index
}
