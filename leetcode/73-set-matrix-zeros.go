package main

import "fmt"

func main() {
	input := [][]int{
		{0, 1, 2, 0},
		{3, 4, 5, 2},
		{1, 3, 1, 5},
	}

	setZeroes(input)
}

func setZeroes(matrix [][]int) {
	var zeroPosition [][]int
	for x, v := range matrix { // x axes
		for y, vv := range v { // y axes
			fmt.Println(x, y, vv)
			if vv == 0 {
				zeroPosition = append(zeroPosition, []int{x, y})
			}
		}
	}

	// changing the x axes
	for _, v := range zeroPosition {
		for i := 0; i < len(matrix[0]); i++ {
			x := v[0]
			matrix[x][i] = 0
		}
	}
	// //changing the y axes
	for _, v := range zeroPosition {
		for i := 0; i < len(matrix); i++ {
			y := v[1]
			matrix[i][y] = 0
		}
	}
}

/*
	1) go over the matrix to find the zeros
	2) When find. Fill the rows and columns with zeros as well
		2.1)

	3x3

	zero in the position 1x1

	1x0 1x1 1x2
	0x1 1x1 2x1



*/
