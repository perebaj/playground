package main

import (
	"fmt"
	"sort"
)

func main() {
	redShirtSpeeds := []int{5, 5, 3, 9, 2}
	blueShirtSpeeds := []int{3, 6, 7, 2, 1}
	fastest := false

	fmt.Println(TandemBicycle(redShirtSpeeds, blueShirtSpeeds, fastest))
}

// The tricky of this problem is: if we compare 2 guys. the guy that have the higher speed will be the one to be select. Idenpendently
// if we are looking to maximize or minimize the "equation"
func TandemBicycle(redShirtSpeeds []int, blueShirtSpeeds []int, fastest bool) int {
	// Write your code here.
	// tandemSpeed = max(speedA, speedB)
	// fastest = true max
	// fastest = false min
	/*
		arr1 = 5,5,3,9,2

		sort = 2,3,5,5,9
		sort = 1,2,3,6,7
		sortAll = 1,2,2,3,3,5,5,6,7,9

		arr2 = 3,6,7,2,1
		5,6 | 1, 2

		5,6

		fastest = true


	*/
	sort.Ints(redShirtSpeeds)
	sort.Ints(blueShirtSpeeds)
	fmt.Printf("redShirtSpeeds %v\t Blue %v\n", redShirtSpeeds, blueShirtSpeeds)
	var total int
	var indexR int
	var indexB int
	if fastest {
		indexR, indexB = len(redShirtSpeeds)-1, len(blueShirtSpeeds)-1
		for i := len(redShirtSpeeds); i > 0; i-- {
			if redShirtSpeeds[indexR] >= blueShirtSpeeds[indexB] {
				fmt.Println(redShirtSpeeds[indexR])
				total = total + redShirtSpeeds[indexR]
				indexR--
			} else {
				fmt.Println(blueShirtSpeeds[indexB])
				total = total + blueShirtSpeeds[indexB]
				indexB--
			}
		}
	} else {
		for i := 0; i < len(redShirtSpeeds); i++ {
			if redShirtSpeeds[i] >= blueShirtSpeeds[i] {
				total = total + redShirtSpeeds[i]
			} else {
				total = total + blueShirtSpeeds[i]
			}
		}
	}

	fmt.Println(total)
	return total
}
