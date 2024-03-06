package main

import (
	"fmt"
	"slices"
)

func main() {
	candies := []int{2, 3, 5, 1, 3}
	extraCandies := 3

	fmt.Println(kidsWithCandies(candies, extraCandies))
}

func kidsWithCandies(candies []int, extraCandies int) []bool {
	max := slices.Max(candies)
	fmt.Println(max)

	result := make([]bool, len(candies))

	for i := 0; i < len(candies); i++ {
		if candies[i]+extraCandies >= max {
			result[i] = true
		} else {
			result[i] = false
		}
	}

	return result
}
