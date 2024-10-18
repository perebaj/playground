package main

import (
	"fmt"
	"math"
	"sort"
)

/*
1. Code Question 1

A user is using the Amazon fitness tracker, and they are engaged in a jumping exercise routine.
The user is positioned on the ground, and there are n stones, each placed at different heights.
 The height of the i-th stone is represented by height[i] meters.

The goal is to maximize the calorie burn during this exercise,
and the calories burned when jumping from the j-th stone to the i-th stone is
calculated as (height[i] - height[j])².

The user intends to practice jumping on each stone exactly once but can choose the
order in which they jump. Since the user is looking to optimize their calorie burn
for this single session, the task is to determine the maximum amount of calories
that can be burned during the exercise.

Formally, given an array height of size n, representing the height of each stone,
find the maximum amount of calories the user can burn.

Note that the user can jump from any stone to any stone, and the ground's height is 0
and once the user jumps to a stone from the ground, he/she can never go back to the ground.

Example

n = 3
height = [5, 2, 5]
The user will jump in the following way:
Ground → 3rd stone → 2nd stone → 1st stone

*/

func main() {
	height := []int{3, 4, 2, 55, 6, 5}
	height = append(height, 0)

	sort.Ints(height) // 0,2,5,5

	index1 := 0
	index2 := len(height) - 1

	var calories float64
	for i := 0; i < len(height)-1; i++ {
		if i%2 == 0 {
			fmt.Println(index1, index2)
			fmt.Println(height[index1], height[index2])
			calories = calories + math.Pow(float64(height[index1])-float64(height[index2]), 2)
			index1++
		} else if i%2 != 0 {
			fmt.Println(index1, index2)
			fmt.Println(height[index1], height[index2])
			calories = calories + math.Pow(float64(height[index1])-float64(height[index2]), 2)
			index2--
		}

		fmt.Println("Calories", calories)
	}
}
