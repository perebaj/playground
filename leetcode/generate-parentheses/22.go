package main

import "fmt"

func main() {
	fmt.Println(generateParenthesis(2))
	fmt.Println(generateParenthesis(1))
}

func generateParenthesis(n int) []string {
	var currentResp string
	var stack []string
	backtrack(0, 0, n, &stack, currentResp)
	return stack
}

func backtrack(open, close, n int, stack *[]string, currentResp string) {
	//zero case
	if len(currentResp) == 2*n {
		// return the result
		*stack = append(*stack, currentResp)

		return
	}

	//add open
	if open < n {
		aux := currentResp + "("
		backtrack(open+1, close, n, stack, aux)
	}

	if close < open {
		aux := currentResp + ")"
		backtrack(open, close+1, n, stack, aux)
	}
}
