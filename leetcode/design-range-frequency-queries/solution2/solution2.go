package main

import (
	"fmt"
	"sort"
)

type RangeFreqQuery struct {
	Frequency map[int][]int // map the value to the index where this values are appearing
}

func Constructor(arr []int) RangeFreqQuery {
	m := make(map[int][]int)

	for index, value := range arr {
		m[value] = append(m[value], index) // eaeach sub array is by default sorted, right? We are just appending index that we are traversing from the lower to the upper
	}

	return RangeFreqQuery{
		Frequency: m,
	}
}

func (r *RangeFreqQuery) Query(left int, right int, value int) int {
	subSlice := r.Frequency[value]
	/*
		Never in my life I would be able to think about this solution by myself
		I read the solutions from leetcode and understand the logic

		But the moral of this problem is the following:
		If we saved all indexes where the value appears, your goal is to find if the left and right bound are inside the subSlice, if
		they aren't, we need to find the closest index that match the condition.

		For this reason we are using the sort.Search, this function basically returns the desired place where the value should be inserted or
		acessed.
	*/
	leftBound := sort.Search(len(subSlice), func(i int) bool { return subSlice[i] >= left })
	rightBound := sort.Search(len(subSlice), func(i int) bool { return subSlice[i] >= right+1 })
	fmt.Println(subSlice)
	fmt.Println(leftBound, rightBound)
	return rightBound - leftBound
}

func main() {
	/*
		Cmplexity analysis:
		Constructor: O(n). Memory O(n)
		Query: 2 * O(log n) approximately O(log n). Memory O(1)
	*/

	array := []int{12, 33, 4, 56, 22, 2, 34, 33, 22, 12, 34, 56}
	n := Constructor(array)
	fmt.Println(n.Frequency)
	fmt.Println(n.Query(1, 2, 4))
	fmt.Println(n.Query(0, 11, 33))
}
