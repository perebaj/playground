package main

import (
	"fmt"
	"sort"
)

func main() {
	nums := []int{1, 1, 2}
	k := removeDuplicates(nums)
	fmt.Println(k)
}

func removeDuplicates(nums []int) int {
	auxIndex := 0
	var k int
	for i := 1; i < len(nums); i++ {
		if nums[auxIndex] == nums[i] {
			nums[i] = 101
			k++
		} else {
			auxIndex = i
		}
	}
	sort.Ints(nums)

	for i := 0; i < len(nums); i++ {
		if nums[i] == 101 {
			nums[i] = 0
		}
	}
	fmt.Println(nums)

	return len(nums) - k
}
