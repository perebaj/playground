package main

import (
	"fmt"
	"math"
)

func main() {
	nums := []int{3, 4, 9, 1, 3, 9, 5}
	key := 9
	k := 1
	fmt.Println(findKDistantIndices(nums, key, k))
}

func findKDistantIndices(nums []int, key int, k int) []int {
	var jsSlice []int // all ocurencies of key in the nums slice (j's)
	for k, v := range nums {
		if v == key {
			jsSlice = append(jsSlice, k)
		}
	}

	fmt.Println(jsSlice)
	var resp []int
	for keyNums, _ := range nums {
		for _, v := range jsSlice {
			if math.Abs(float64(keyNums-v)) <= float64(k) {
				resp = append(resp, keyNums)
				break
			}
		}
	}

	return resp
}
