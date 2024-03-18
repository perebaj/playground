package main

import "fmt"

func main() {
	nums := []int{1, 2, 3, 5}
	fmt.Println(containsDuplicate(nums))
}

func containsDuplicate(nums []int) bool {
	m := make(map[int]int)
	for i := 0; i < len(nums); i++ {
		_, ok := m[nums[i]]
		if ok {
			return true
		} else {
			m[nums[i]]++
		}
	}
	return false
}

/*
	O(n)
*/
