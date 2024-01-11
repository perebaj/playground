package main

import "fmt"

func main() {
	prices := []int{1, 2, 4}
	fmt.Println(maxProfit(prices))
}

func maxProfit(prices []int) int {
	if len(prices) < 2 {
		return 0
	}

	min := prices[0]
	// profit := prices[1] - min

	if len(prices) == 2 {
		if min < prices[1] {
			return prices[1] - min
		}
	}

	var profit int
	for i := 0; i < len(prices); i++ {
		if i == len(prices)-1 {
			if profit < prices[i]-min {
				return prices[i] - min
			}
			break
		}
		if min > prices[i+1] {
			min = prices[i+1]

		} else {
			profit = prices[i] - min
		}
	}
	return profit
}
