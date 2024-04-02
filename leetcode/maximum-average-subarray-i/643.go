package main

import (
	"fmt"
	"math"
)

func main() {
	nums := []int{1, 12, -5, -6, 50, 3}
	k := 4
	fmt.Println(findMaxAverage(nums, k))
}

func findMaxAverage(nums []int, k int) float64 {
	startPointer := 0
	endPointer := k - 1
	mean := math.Inf(-1)

	if len(nums) == 1 {
		return float64(nums[0])
	}

	s := sum(nums[startPointer : endPointer+1])
	for {
		auxMean := s / float64(k)
		if auxMean > mean {
			mean = auxMean
		}
		firstElement := nums[startPointer]
		startPointer++
		endPointer++
		if endPointer <= len(nums)-1 {
			s = s - float64(firstElement) + float64(nums[endPointer])
		} else {
			break
		}
	}
	return mean
}

func sum(nums []int) float64 {
	var sum float64
	for _, v := range nums {
		sum += float64(v)
	}
	return sum
}
