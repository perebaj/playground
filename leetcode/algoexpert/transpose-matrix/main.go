package main

import "fmt"

func main() {
	matrix := [][]int{
		{1, 2},
	}

	fmt.Println(len(matrix))    // len linhas
	fmt.Println(len(matrix[0])) // len colunas
	// result := make([][]int, len(matrix))
	var result [][]int
	for col := 0; col < len(matrix[0]); col++ {
		var auxArr []int
		for lin := 0; lin < len(matrix); lin++ {
			fmt.Print(matrix[lin][col])
			auxArr = append(auxArr, matrix[lin][col])
		}
		fmt.Println("")
		result = append(result, auxArr)
	}

	fmt.Println(result)
}
