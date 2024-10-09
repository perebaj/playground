package main

import "fmt"

func SpiralTraverse(array [][]int) []int {
	// Write your code here.
	var result []int
	if len(array) < 2 {
		result = append(result, array[0][0])
		return result
	}
	startRow := 0
	endRow := len(array) - 1
	startCol := 0
	endCol := len(array[0]) - 1
	for startRow <= endRow && startCol <= endCol {
		//main for loop <- im not sure about this one
		// horizontal left to right
		for i := startCol; i <= endCol; i++ {
			//do something
			result = append(result, array[startRow][i])
		}
		startRow++
		//vertical up to down
		for k := startRow; k <= endRow; k++ {
			result = append(result, array[k][endCol])
		}
		endCol--
		// horizontal right to left
		for j := endCol; j >= startCol; j-- {
			result = append(result, array[endRow][j])
			//do something that i dont know, yet!
		}
		endRow--
		//vertical down to up

		for h := endRow; h >= startRow; h-- {
			result = append(result, array[h][startCol])
			//do something
		}
		startCol++
		fmt.Println(result)
	}
	return result

}

func SpiralTraverse2(array [][]int) []int {
	spiral := make([]int, 0)
	total := len(array) * len(array[0])
	xMin, xMax, yMin, yMax := 0, len(array)-1, 0, len(array[0])-1

	for step := 0; step < total; {
		// Traverse from left to right across the top row
		for i := yMin; i <= yMax && step < total; i++ {
			spiral = append(spiral, array[xMin][i])
			step++
		}
		xMin++

		// Traverse from top to bottom down the right column
		for i := xMin; i <= xMax && step < total; i++ {
			spiral = append(spiral, array[i][yMax])
			step++
		}
		yMax--

		// Traverse from right to left across the bottom row
		for i := yMax; i >= yMin && step < total; i-- {
			spiral = append(spiral, array[xMax][i])
			step++
		}
		xMax--

		// Traverse from bottom to top up the left column
		for i := xMax; i >= xMin && step < total; i-- {
			spiral = append(spiral, array[i][yMin])
			step++
		}
		yMin++
	}

	return spiral
}
