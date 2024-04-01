package main

import (
	"fmt"
	"math"
	"sort"
)

func main() {
	input := []int{9, 4, 1, 7}
	k := 2
	fmt.Println(minimumDifference(input, k))

}

func minimumDifference(nums []int, k int) int {
	sort.Ints(nums)
	fmt.Println(nums)
	min := math.Inf(1)
	pointerStart := 0
	pointerEnd := k - 1
	for pointerEnd < len(nums) {
		fmt.Println(nums[pointerEnd], nums[pointerStart])
		aux := float64(nums[pointerEnd] - nums[pointerStart])
		if aux < min {
			min = aux
		}
		pointerStart++
		pointerEnd++
	}

	return int(min)
}
