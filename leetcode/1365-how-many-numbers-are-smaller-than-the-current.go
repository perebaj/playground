package main

import "fmt"

func main() {
	nums := []int{8, 1, 2, 2, 3}
	fmt.Println(smallerNumbersThanCurrent(nums))
}

func smallerNumbersThanCurrent(nums []int) []int {
	aux := make([]int, len(nums))
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			// fmt.Println(i, j)
			if nums[i] == nums[j] {
				continue
			}
			if nums[i] > nums[j] {
				aux[i]++
			} else {
				aux[j]++
			}
		}
	}

	return aux
}
