package main

import "fmt"

func main() {
	fmt.Println(climbStairs(5))
	// fmt.Println(solve(3))
}

func climbStairs(n int) int {
	s := make([]int, n+1)
	s[0] = 1
	s[1] = 1
	for i := 2; i <= n; i++ {
		s[i] = s[i-1] + s[i-2]
	}
	return s[n]
}
