package main

import "fmt"

func main() {
	nums := []int{1, 0, 1}
	moveZeroes(nums)
}

func moveZeroes(nums []int) {
	var idx1, idx2 int
	for {
		if nums[idx1] == 0 && nums[idx2] != 0 {
			nums[idx1] = nums[idx2]
			nums[idx2] = 0
			idx1++
		}
		if nums[idx1] != 0 {
			idx1++
		}
		idx2++
		if idx2 >= len(nums) {
			break
		}
	}
	fmt.Println(nums)
}

/*
	0,0 = 1,1
	1,1 = 0,0






	2 pointers
		first pointer to indicate the current index when was zero
		second pointer to find for a non-zero value
	idx1, idx2
	0,      0
	nums[i1], nums[i2] = 0, 0

	0,      1 = 0,1
	[1,0,0,3,12]

	1       2 = 0,0

	1 3, | 0, 3
	[1,3,0,0,12]


*/
