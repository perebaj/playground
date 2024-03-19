package main

import (
	"fmt"
	"reflect"
	"slices"
)

func main() {
	mainValue := []int{49, 18, 16}
	sub := [][]int{
		{16, 18, 49},
	}
	fmt.Println(solution(mainValue, sub))

	mainValue = []int{15, 88, 88}
	sub = [][]int{
		{88},
		{15},
		{88},
	}
	fmt.Println(solution(mainValue, sub))
}

func solution(mainValue []int, sub [][]int) bool {
	for _, v := range sub {
		if len(v) > 1 {
			b, start, end := containsSubArray(mainValue, v)
			if b {
				mainValue = slices.Delete(mainValue, start, end)
			}
		} else {
			ok, index := contains(mainValue, v[0])
			if ok {
				mainValue = slices.Delete(mainValue, index, index+1)
			}
		}
	}
	return len(mainValue) == 0
}

func contains(s []int, target int) (bool, int) {
	for k, v := range s {
		if v == target {
			return true, k
		}
	}
	return false, -1
}

func containsSubArray(mainSlice []int, subSlice []int) (bool, int, int) {
	start := 0
	end := len(subSlice)
	// 0/
	for end <= len(mainSlice) {
		if !reflect.DeepEqual(mainSlice[start:end], subSlice) {
			start++
			end++
		} else {
			return true, start, end
		}
	}

	return false, start, end
}
