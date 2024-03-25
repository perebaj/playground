package main

import "fmt"

func main() {
	nums := []int{1, 2, 2, 2, 2, 2, 2, 4, 4, 5, 6}
	fmt.Println(binarySearch(nums, 2))

	fmt.Println(binarySearchLowerBound(nums, 2)) // 1
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

// return the first ocurrence where target should be put
func binarySearchLowerBound(nums []int, target int) int {
	left := 0
	right := len(nums)
	for left < right {
		middle := (left + right) / 2
		if nums[middle] < target { // true value
			left = middle + 1
		} else {
			// left = middle + 1
			right = middle
		}
	}

	return left
}

func binarySearchUpperBound(nums []int, target int) int {
	left := 0
	right := len(nums)
	for left < right {
		middle := (left + right) / 2
		if nums[middle] <= target { // true value
			left = middle + 1
		} else {
			// left = middle + 1
			right = middle
		}
	}

	return left
}
