package main

import "fmt"

func main() {
	accounts := [][]int{
		{1, 5},
		{7, 3},
		{3, 5},
	}

	fmt.Println(maximumWealth(accounts))
}

func maximumWealth(accounts [][]int) int {
	var max int
	for _, v := range accounts {
		var aux int
		for _, vv := range v {
			aux = aux + vv
		}

		if aux > max {
			max = aux
		}
	}
	return max
}
