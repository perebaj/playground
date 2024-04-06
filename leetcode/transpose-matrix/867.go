package main

import "fmt"

func main() {
	matrix := [][]int{
		{2, 4, -1},
		{-10, 5, 11},
		{18, -7, 6},
	}

	matrix = [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	fmt.Println(transpose(matrix))
}

func transpose(matrix [][]int) [][]int {
	colSize := len(matrix)
	rowSize := len(matrix[0])

	newMatrix := make([][]int, rowSize)
	for i := 0; i < rowSize; i++ {
		newMatrix[i] = make([]int, colSize)
	}
	fmt.Println(newMatrix)

	for j := 0; j < colSize; j++ {
		for i := 0; i < rowSize; i++ {
			newMatrix[i][j] = matrix[j][i]
		}
	}
	return newMatrix
}
