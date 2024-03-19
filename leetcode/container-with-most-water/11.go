package main

import "fmt"

func main() {
	height := []int{1, 1}
	fmt.Println(maxArea(height))
}

func maxArea(height []int) int {
	start := 0
	end := len(height) - 1
	result := 0
	curArea := 0
	for start < end {
		curArea = (end - start) * min(height[start], height[end])
		if curArea > result {
			result = curArea
		}

		if height[start] <= height[end] {
			start++
		} else {
			end--
		}
	}
	return result
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
