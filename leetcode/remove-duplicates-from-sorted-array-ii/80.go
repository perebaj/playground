package main

import (
	"fmt"
)

type Person struct {
	Name string
	Age  string
}

func main() {
	nums := []int{0, 0, 1, 1, 1, 1, 2, 3, 3}
	fmt.Println(removeDuplicates(nums))
}

/*
I take hours thinking in this problem, but when I finish to solve I beat 100% of the users, using my first idea.

Very happy day
*/
func removeDuplicates(nums []int) int {
	m := make(map[int]int)
	var position int
	var result int
	for i := 1; i < len(nums); i++ {
		_, ok := m[nums[position]]
		if !ok {
			m[nums[i]]++
		}
		if m[nums[i]] == m[nums[position]] && m[nums[position]] < 2 {
			position++
			m[nums[i]]++
			nums[position] = nums[i]
			result++
		} else if m[nums[i]] != m[nums[position]] {
			position++
			nums[position] = nums[i]
			result++
		}
	}

	return result + 1
}
