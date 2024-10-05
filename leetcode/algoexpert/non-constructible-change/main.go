package main

import (
	"sort"
)

func NonConstructibleChange(coins []int) int {
	sort.Ints(coins)

	minChange := 0

	for _, coin := range coins {
		if coin > minChange+1 {
			break
		}
		minChange += coin
	}

	return minChange + 1
}

func main() {
	res := NonConstructibleChange([]int{5, 7, 1, 1, 2, 3, 22})
	println(res)
}
