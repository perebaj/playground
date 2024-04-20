// https://www.algoexpert.io/questions/tournament-winner
package main

import "fmt"

func TournamentWinner(competitions [][]string, results []int) string {
	// Write your code here.
	m := make(map[string]int)
	for k, v := range competitions {
		aux := results[k]
		if aux == 0 {
			aux = 1
		} else if aux == 1 {
			aux = 0
		}
		m[v[aux]] += 3
	}
	fmt.Println(m)

	var result string
	var max int
	for k, v := range m {
		if v > max {
			max = v
			result = k
		}
		fmt.Println(result, max)
	}

	return result
}

// 1 home team won
// 0 away team won
/*
   {
       HTML: Pontuation
       C#: Pontuation
       Python: Pontuation
   }
*/
