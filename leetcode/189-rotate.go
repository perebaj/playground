package main

import "fmt"

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7}
	rotate(nums, 3)
}

func rotate(nums []int, k int) {
	numsAux := make([]int, len(nums))
	for i := 0; i < k; i++ {
		for j := 0; j < len(nums); j++ {
			if j == len(nums)-1 {
				numsAux[0] = nums[j]
			} else {
				numsAux[j+1] = nums[j]
			}
		}
		copy(nums, numsAux)
		fmt.Println("copy:", numsAux)

	}
	fmt.Println("nums aux:", numsAux)
}
