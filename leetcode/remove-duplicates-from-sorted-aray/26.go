package main

import "fmt"

func main() {
	// fmt.Println(removeDuplicates(nums))
	nums := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	fmt.Println(removeDuplicates2(nums))
	nums = []int{1, 1, 2}
	fmt.Println(removeDuplicates2(nums))
}

func removeDuplicates(nums []int) int {
	m := make(map[int]int)
	var total int
	position := 0
	for i := 0; i < len(nums); i++ {
		_, ok := m[nums[i]]
		if !ok {
			m[nums[i]] = 0
			nums[position] = nums[i]
			position++
			total++
		}
	}
	fmt.Println(nums)

	return total
}

/*
Two pointer approach is the better way to solve it
one pointer will point to the last unic elemet that was inserte, and the other one will walk through the array, trying to find different elements
*/
func removeDuplicates2(nums []int) int {
	// var total int
	position := 0
	var total int
	for i := 1; i < len(nums); i++ {
		if nums[i] != nums[position] && nums[i] > nums[position] {
			fmt.Println(nums[position], nums[i])
			position++
			nums[position] = nums[i]
			total++
		}
	}
	// fmt.Println(nums)
	return total + 1
}
