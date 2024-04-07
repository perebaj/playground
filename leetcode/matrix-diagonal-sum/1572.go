package main

import "fmt"

func main() {
	mat := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	fmt.Println(mat)
}

func diagonalSum(mat [][]int) int {
	matSize := len(mat)

	//main diagonal: i == j
	//secondary diagonal: i+j=matSize
	if matSize == 1 {
		return mat[0][0]
	}

	mainDiagonal := 0
	secondaryDiagonal := 0
	for i := 0; i < matSize; i++ {
		for j := 0; j < matSize; j++ {
			if i == j {
				mainDiagonal += mat[i][j]
			} else if i+j == matSize-1 {
				secondaryDiagonal += mat[i][j]
			}
		}
	}

	return secondaryDiagonal + mainDiagonal
}
