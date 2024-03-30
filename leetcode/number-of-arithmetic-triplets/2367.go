package main

import "fmt"

func main() {
	nums := []int{0, 1, 4, 6, 7, 10}
	diff := 3
	fmt.Println(arithmeticTriplets(nums, diff)) // 2
	nums = []int{4, 5, 6, 7, 8, 9}
	diff = 2
	fmt.Println(arithmeticTriplets(nums, diff)) // 2
}

func arithmeticTriplets(nums []int, diff int) int {
	m := make(map[int]int)
	for k, v := range nums {
		m[v] = k
	}
	var resp int
	fmt.Println(m)
	for j, _ := range m {
		_, ok := m[j-diff]
		_, ok2 := m[diff+j]
		fmt.Println(j-diff, j, diff+j, ok, ok2)
		if ok && ok2 && (j-diff) < j && j < (diff+j) {
			resp++
		}
	}
	return resp
}
