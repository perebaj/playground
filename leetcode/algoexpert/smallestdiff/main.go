package main

import (
	"fmt"
	"math"
	"sort"
)

func main() {
	arrayOne := []int{-1, 5, 10, 20, 28, 3}
	arrayTwo := []int{26, 134, 135, 15, 17}

	res := SmallestDifference(arrayOne, arrayTwo)
	fmt.Println(res)
}

func SmallestDifference(array1, array2 []int) []int {
	// Write your code here.
	sort.Ints(array1)
	sort.Ints(array2)

	var indexArr1 int
	var indexArr2 int
	lenArr1 := len(array1)
	lenArr2 := len(array2)

	diff := math.MaxInt32
	var res []int
	for indexArr1 < lenArr1 && indexArr2 < lenArr2 {
		val1, val2 := array1[indexArr1], array2[indexArr2]

		diffAux := absDiffInt(val1, val2)
		if diffAux < diff {
			res = []int{val1, val2}
			diff = diffAux
		}

		if val1 < val2 {
			indexArr1++
		} else {
			indexArr2++
		}
	}
	return res
}

func absDiffInt(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}
