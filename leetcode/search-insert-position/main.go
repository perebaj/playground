package main

import "fmt"

func searchInsert(nums []int, target int) int {
	return binarySearch(nums, target)

}

func binarySearch(nums []int, target int) int {
	left := 0
	right := len(nums)

	for left < right {
		middle := (right + left) / 2
		if nums[middle] < target {
			left = middle + 1
		} else {
			right = middle
		}
	}

	return left
}

func main() {
	nums := []int{1, 3, 5, 6}
	fmt.Println(searchInsert(nums, 5))
}
