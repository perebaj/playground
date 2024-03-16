package main

import "fmt"

func main() {
	nums := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	fmt.Println(removeDuplicates(nums))
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
