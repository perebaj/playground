package main

import (
	"fmt"
)

func main() {
	nums := []int{5, 3, 1, 1, 1, 3, 5, 73, 1}
	fmt.Println(topKFrequent(nums, 3))
}

func topKFrequent(nums []int, k int) []int {
	m := make(map[int]int)
	for _, v := range nums {
		m[v]++
	}

	fmt.Println(m)

	m2 := make(map[int][]int)
	for k, v := range m {
		_, ok := m2[v]
		if !ok {
			m2[v] = []int{k}
		} else {
			m2[v] = append(m2[v], k)
		}
	}
	fmt.Println(m2)
	var result []int
	for i := len(nums); i > 0; i-- {
		_, ok := m2[i]
		if ok && len(result) < k {
			values := m2[i]
			fmt.Println(values, i)
			result = append(result, values...)

		}

	}
	fmt.Println(result)
	return result
}
