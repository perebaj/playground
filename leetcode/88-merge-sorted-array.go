package main

import (
	"fmt"
	"sort"
)

func main() {
	nums1 := []int{1, 2, 3, 0, 0, 0}
	nums2 := []int{2, 5, 6}

	merge(nums1, 3, nums2, 3)
}

func merge(nums1 []int, m int, nums2 []int, n int) {
	// append nums2 to nums1
	for i := 0; i < n; i++ {
		nums1[m+i] = nums2[i]
	}
	fmt.Printf("nums1: %v\n", nums1)
	sort.Ints(nums1)
	fmt.Printf("nums1: %v\n", nums1)
}
