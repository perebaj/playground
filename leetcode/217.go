package main

import "fmt"

func main() {
	in := []int{1, 2, 3, 1}
	fmt.Println(containsDuplicate(in))

}

func containsDuplicate(nums []int) bool {
	m := make(map[int]int)

	for i := 0; i < len(nums); i++ {
		m[nums[i]]++
	}

	for _, v := range m {
		if v > 1 {
			return true
		}
	}
	return false
}
