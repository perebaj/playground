package main

import (
	"sort"
)

func MinimumWaitingTime(queries []int) int {
	sort.Ints(queries)
	var total int
	index := 1
	for i := 0; i < len(queries)-1; i++ {
		total = total + (queries[i] * (len(queries) - index))
		index++
	}
	return total
}
