package main

import (
	"fmt"
	"sort"
)

func main() {
	nums := []int{0, 0, 0, 0}
	fmt.Println(threeSum(nums))
}

/*
This one is a hard problem, even saying that is a medium, the tricky in the end, doing an infinit for loop to avoid conflicts/duplicates result
was the hardest part, and also the sorting

I will try to review it in a near future to fix it on my head
*/
func threeSum(nums []int) [][]int {
	sort.Ints(nums)
	fmt.Println(nums)
	var result [][]int
	for i, v := range nums {
		if i > 0 && nums[i-1] == nums[i] {
			continue
		}
		left := i + 1
		right := len(nums) - 1
		for left < right {

			count := v + nums[right] + nums[left]
			if (count) > 0 { // means that we need to decrease the  right pointer{
				right--
			} else if count < 0 {
				left++
			} else {
				result = append(result, []int{v, nums[left], nums[right]})
				// if nums[left] == nums[left-1] {
				// 	left = left + 2
				// } else {
				// 	left++S
				// }
				left++
				for nums[left] == nums[left-1] && left < right {
					left++
				}

			}

		}

	}

	return result
}
