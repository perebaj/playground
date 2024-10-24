package main

import "sort"

func ClassPhotos(redShirtHeights []int, blueShirtHeights []int) bool {
	// Write your code here.
	/*
		[5, 8, 1, 3, 4],
		[2, 9, 6, 4, 5]
		1,3,4,5,8
		2,4,5,6,9
	*/
	sort.Ints(redShirtHeights)
	sort.Ints(blueShirtHeights)

	isRedFrontLine := false
	isBlueFrontLine := false

	for i := 0; i < len(redShirtHeights); i++ {
		if redShirtHeights[i] > blueShirtHeights[i] {
			isBlueFrontLine = true // they are smaller than the red line
		} else if redShirtHeights[i] < blueShirtHeights[i] {
			isRedFrontLine = true
		} else {
			return false // they are equal, this is not a valid rule
		}

		if isBlueFrontLine && isRedFrontLine {
			return false
		}
	}
	return true
}
