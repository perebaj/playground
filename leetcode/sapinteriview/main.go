package main

import "fmt"

func main() {
	matrix := [][]int{
		{7},
		{9},
		{6},
	}
	fmt.Println(spiralOrder(matrix))
}

// leetcode 54
func spiralOrder(matrix [][]int) []int {
	xMax := len(matrix) - 1
	xMin := 0
	yMax := len(matrix[len(matrix)-1]) - 1
	yMin := 0

	var resp []int
	if len(matrix[0]) < 2 {
		for _, v := range matrix {
			resp = append(resp, v[0])
		}
		return resp
	}

	// finalLen
	for {
		for j := yMin; j <= yMax; j++ {
			resp = append(resp, matrix[xMin][j])
		}
		if len(resp) == len(matrix)*len(matrix[0]) {
			break
		}
		xMin++
		for i := xMin; i <= xMax; i++ {
			resp = append(resp, matrix[i][yMax])
		}
		if len(resp) == len(matrix)*len(matrix[0]) {
			break
		}
		yMax--
		for j := yMax; j >= yMin; j-- {
			resp = append(resp, matrix[xMax][j])
		}
		if len(resp) == len(matrix)*len(matrix[0]) {
			break
		}
		xMax--
		for i := xMax; i >= xMin; i-- {
			resp = append(resp, matrix[i][yMin])
		}
		if len(resp) == len(matrix)*len(matrix[0]) {
			break
		}
		yMin++
		fmt.Println(resp)
	}
	return resp
}
