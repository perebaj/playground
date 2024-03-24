package main

import (
	"fmt"
	"slices"
	"sort"
)

// Important aspects of the sort.Search function:
// The index return the desired position of the target value,
// if you verify that array[index] == target, then the target exists in
// the array, otherwise, the target does not exist but the index is
// the right place to insert the target value.
func sortSearch() {
	array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	array2 := []int{9, 8, 7, 6, 5, 4, 3, 2, 1}
	target := 1
	index := sort.Search(len(array2), func(i int) bool {
		return array2[i] <= target
	})

	index2 := sort.Search(len(array), func(i int) bool {
		return array[i] >= target
	})

	fmt.Println(index)
	fmt.Println(index2)

	fmt.Println(array2[index])
	fmt.Println(array[index2])
}

func Index() int {
	array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	index := slices.Index(array, 1) // 0

	return index

}

func Constructor(arr []int) {
	indexMap := make(map[int][]int)
	for i, num := range arr {
		indexMap[num] = append(indexMap[num], i)
	}

	fmt.Println(indexMap)
}
