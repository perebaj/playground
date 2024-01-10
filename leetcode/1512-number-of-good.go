package main

import "fmt"

func main() {
	nums := []int{1, 2, 3}

	fmt.Println(numIdenticalPairs(nums))
}

func numIdenticalPairs(nums []int) int {
	var aux int
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i] == nums[j] {
				aux++
			}
		}
	}
	return aux
}
