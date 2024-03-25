package main

import "fmt"

func main() {
	nums := []int{4, 5, 1, 2, 3}
	fmt.Println(findMin(nums))
}

func findMin(nums []int) int {
	left := 0
	right := len(nums) - 1
	var small int

	small = nums[right]
	for left <= right {
		if nums[right] < nums[left] { // messy case
			middle := (left + right) / 2
			if middle < len(nums) && nums[middle] < small {
				small = nums[middle]
				right = middle + 1
			} else {
				left = middle + 1
			}
		} else if right > left {
			return nums[left]
		} else {
			break
		}
	}
	return small
}
