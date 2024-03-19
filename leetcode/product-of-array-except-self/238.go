package main

import "fmt"

func main() {
	nums := []int{1, 2, 3, 4}
	fmt.Println(productExceptSelf(nums))
}

func productExceptSelf(nums []int) []int {
	left := make([]int, len(nums))
	right := make([]int, len(nums))
	left[0] = 1
	right[len(nums)-1] = 1
	for i := 1; i < len(nums); i++ {
		left[i] = left[i-1] * nums[i-1]
	}

	result := make([]int, len(nums))

	for i := len(nums) - 2; i >= 0; i-- {
		right[i] = right[i+1] * nums[i+1]
		fmt.Println(right[i+1], left[i+1])
		result[i+1] = right[i+1] * left[i+1]
	}

	// fmt.Println(result)
	result[0] = right[0] * left[0]
	// var result []int
	// for i := 0; i < len(nums); i++ {
	// 	result = append(result, left[i]*right[i])
	// }

	return result
}
