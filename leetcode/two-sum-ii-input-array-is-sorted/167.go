package main

import (
	"fmt"
	"sort"
)

func main() {
	numbers := []int{2, 7, 11, 15}
	target := 9
	fmt.Println(twoSumTwoPointers(numbers, target))

	numbers = []int{0, 0, 3, 4}
	target = 0
	fmt.Println(twoSumTwoPointers(numbers, target))

}

func TwoSum(numbers []int, target int) []int {
	for i := 0; i < len(numbers); i++ {
		aux := target - numbers[i]
		idx2 := sort.Search(len(numbers), func(j int) bool { return numbers[j] >= aux && j > i })
		if idx2 < len(numbers) && numbers[idx2] == aux {
			return []int{i + 1, idx2 + 1}
		}
	}
	return nil
}

func twoSumTwoPointers(numbers []int, target int) []int {
	pointer1 := 0
	pointer2 := len(numbers) - 1
	for {
		if numbers[pointer1]+numbers[pointer2] == target {
			break
		}
		if numbers[pointer1]+numbers[pointer2] > target {
			pointer2--
		} else {
			pointer1++
		}

	}

	// return []int{pointer1 + 1, pointer2 + 1}
	return []int{pointer1 + 1, pointer2 + 1}
}

/*
O(NlogN)

we are iterating over the numbers array N times and searching N times as well. Preatty easy to find this Big O answer

I struggled a little bit to understand the sort.Search. But was nice
*/
