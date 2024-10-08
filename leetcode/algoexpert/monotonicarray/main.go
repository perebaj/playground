package main

import "fmt"

func IsMonotonic(array []int) bool {

	if len(array) < 2 {
		return true
	}

	// Write your code here.
	decrease := false
	increase := false
	index1 := 0
	index2 := 1

	for index2 < len(array) {
		if array[index2] < array[index1] {
			decrease = true
		} else if array[index2] > array[index1] {
			increase = true
		}

		fmt.Println(decrease, increase)
		if increase && decrease {
			return false
		}

		index1++
		index2++
	}

	return true
}
