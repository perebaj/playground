package main

import "fmt"

func main() {
	in := []int{2, 2, 3, 4}
	fmt.Println(findLucky(in))
}

func findLucky(arr []int) int {
	m := make(map[int]int)

	for _, v := range arr {
		m[v]++
	}
	fmt.Println(m)
	var auxMax int

	for k, v := range m {
		fmt.Println(k, v)
		if k == v && v > auxMax {
			auxMax = v
		}
	}
	if auxMax > 0 {
		return auxMax
	}
	return -1
}
