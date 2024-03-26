package main

import "fmt"

func searchMatrix(matrix [][]int, target int) bool {
	for _, v := range matrix {
		if target >= v[0] && target <= v[len(v)-1] {
			if target == v[0] || target == v[len(v)-1] {
				return true
			} else {
				resp := binarySeach(v, target)
				fmt.Println(resp)
				return resp
			}
		}
	}
	return false
}

func binarySeach(nums []int, target int) bool {
	left := 0
	right := len(nums) - 1
	var resp bool
	for left <= right {
		middle := (right + left) / 2
		if nums[middle] > target {
			right = middle - 1
		} else if nums[middle] < target {
			left = middle + 1
		} else {
			return true
		}
	}

	return resp
}

func main() {
	matrix := [][]int{
		{1, 3, 5, 7},
		{10, 11, 16, 20},
		{23, 30, 34, 60},
	}
	target := 3
	fmt.Println(searchMatrix(matrix, target))
}
