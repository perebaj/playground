package main

import "fmt"

func main() {
	mat := [][]int{
		{0, 1},
		{1, 0},
	}

	fmt.Println(rowAndMaximumOnes(mat))
}

func rowAndMaximumOnes(mat [][]int) []int {
	result := []int{0, 0} // index of the row, and how many items it has

	for rowIndex, rows := range mat {
		var count int
		for _, element := range rows {
			if element == 1 {
				count++
			}
		}
		if count >= 1 && count > result[1] {
			result[0] = rowIndex
			result[1] = count
		}
	}
	return result
}
