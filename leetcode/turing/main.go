package main

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
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

	min := 15
	max := 25

	fmt.Println(solution2(min, max))
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

func solution2(m int, n int) int {
	mapp := make(map[int]int)
	for i := m; i <= n; i++ {
		numStr := strconv.Itoa(i)
		sumN := Summ(numStr)
		mapp[sumN]++
	}

	fmt.Println(mapp)
	var sSlice []int
	for _, v := range mapp {
		sSlice = append(sSlice, v)
	}

	resp := slices.Max(sSlice)

	return resp
}

func Summ(numStr string) int {
	var sum int
	for _, v := range numStr {
		value, _ := strconv.Atoi(string(v))
		sum = sum + value
	}

	return sum
}
