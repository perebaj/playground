package main

import "fmt"

func main() {
	nums := []int{0, 3, 7, 2, 5, 8, 4, 6, 0, 1}
	fmt.Println(longestConsecutive2(nums))
}

func longestConsecutive(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	m := make(map[int]int)
	for _, v := range nums {
		m[v]++
	}
	var initalElementSlice []int
	for k := range m {
		initial := k - 1
		_, ok := m[initial]
		if !ok {
			initalElementSlice = append(initalElementSlice, k)
		}
	}
	positionInitialElementSlice := 0
	var result int
	curResult := 1
	for {
		_, ok := m[initalElementSlice[positionInitialElementSlice]+curResult]
		if !ok {
			if curResult > result {
				result = curResult
			}
			curResult = 1
			positionInitialElementSlice++
		} else {
			curResult++
		}

		if positionInitialElementSlice >= len(initalElementSlice) {
			break
		}
	}

	return result
}

func longestConsecutive2(nums []int) int {
	m := make(map[int]bool)

	for _, v := range nums {
		m[v] = false
	}

	fmt.Println(m)
	var max int
	for k, _ := range m {
		var aux int
		for i := 0; i <= len(m)+1; i++ {
			_, ok := m[k+i]
			fmt.Println("OK:", ok)
			fmt.Println("Val: and index", k, i)
			if !ok {
				if aux > max {
					max = aux
					fmt.Println("MAX:", max)
				}
				break
			}
			aux++
		}
	}
	return max
}
