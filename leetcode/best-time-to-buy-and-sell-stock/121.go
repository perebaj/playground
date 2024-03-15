package main

import "fmt"

func main() {
	prices := []int{3, 4, 8, 2, 10}
	fmt.Println(maxProfit(prices))
}

func maxProfit(prices []int) int {
	// profit = sell - buy
	buy := 0
	sell := 1
	var maxProfit int
	for sell < len(prices) {
		curProf := prices[sell] - prices[buy]
		if curProf <= 0 {
			buy = sell
		} else if curProf > maxProfit {
			maxProfit = curProf
		}

		sell++
	}
	return maxProfit
}

/*
	in this exercise we need to adopt a 2 pointer approach, trying to find the maxProfit meanwhile we are iterating over prices array.


	The first thing to improve the exercise resolution, is to think about the use cases that we have:
	for example
	when the value of sell its smaller than buy, in that case we have a negative number. In other words. The current sell should be a buy day

	the other cases is just we comparing a fixed variable (maxProfit) with the current profit, and trying to find wich one is bigger

	I thought that this should be a medium leetcode problems instead of easy
*/
