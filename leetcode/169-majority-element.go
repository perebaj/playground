package main

import "fmt"

func main() {
	nums := []int{3, 2, 3}
	k := majorityElement(nums)
	fmt.Println(k)
}

func majorityElement(nums []int) int {
	m := make(map[int]int)
	for _, v := range nums {
		m[v] = 0
	}

	for k := range m {
		var aux int
		for i := 0; i < len(nums); i++ {
			if k == nums[i] {
				aux++
			}
		}
		m[k] = aux
	}

	majElement := len(nums) / 2
	var resp int
	for k, v := range m {
		if v > majElement {
			resp = k
		}
	}
	fmt.Println(majElement, resp)
	return resp
}
