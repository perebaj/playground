package main

import (
	"fmt"
	"sort"
)

func main() {
	in := []int{4, 1, -1, 2, -1, 2, 3}
	k := 2

	fmt.Println(topKFrequent(in, k))
}

func topKFrequent(nums []int, k int) []int {
	m := make(map[int]int)

	for _, v := range nums {
		m[v]++
	}

	// unstructure the count value into a slice
	var sSlice []int
	for _, v := range m {
		sSlice = append(sSlice, v)
	}
	sort.Ints(sSlice)

	aux := sSlice[len(sSlice)-k:]
	fmt.Println(aux)
	fmt.Println(m)
	var resp []int

	for kk, vv := range m {
		for _, vvv := range aux {
			if vvv == vv {
				resp = append(resp, kk)
				break
			}
		}
	}
	return resp
}
