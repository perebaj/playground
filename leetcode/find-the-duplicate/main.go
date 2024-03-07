/*
Find The Duplicates
Given two sorted arrays arr1 and arr2 of passport numbers, implement a function
findDuplicates that returns an array of all passport numbers that are both in
arr1 and arr2. Note that the output array should be sorted in an ascending order.

Let N and M be the lengths of arr1 and arr2, respectively. Solve for two cases
and analyze the time & space complexities of your solutions: M ≈ N - the array
lengths are approximately the same M ≫ N - arr2 is much bigger than arr1.
*/

package main

import (
	"fmt"
)

func FindDuplicates(arr1 []int, arr2 []int) []int {
	// your code goes here
	var result []int
	for _, v := range arr1 {
		_, ok := BinarySearch(arr2, v)
		if ok {
			result = append(result, v)
		}
	}
	return result
}

func main() {
	arr1 := []int{1, 2, 3, 5, 6, 7}
	arr2 := []int{3, 6, 7, 8, 20}

	// fmt.Println(BinarySearch(arr1, 20))
	fmt.Println(FindDuplicates(arr1, arr2))
}

func BinarySearch(S []int, E int) (int, bool) {
	low := 0
	high := len(S) - 1
	// mid := low + (high-low)/2

	for high > low {
		mid := low + (high-low)/2
		if E == S[mid] {
			return S[mid], true
		} else if E < S[mid] {
			high--
		} else {
			low++
		}
	}
	return -1, false
}

/*

 */
