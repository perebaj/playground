package main

import "fmt"

func main() {
	s := "x"
	fmt.Println(countGoodSubstrings(s))
}

func countGoodSubstrings(s string) int {
	var result int
	if len(s) == 1 {
		return 0
	}
	maxSubstrings := len(s) - 3 + 1
	for i := 0; i < maxSubstrings; i++ {
		aux := s[i : i+3]
		m := make(map[rune]int)
		for _, v := range aux {
			_, ok := m[v]
			if ok {
				result++
				break
			} else {
				m[v]++
			}
		}
	}

	fmt.Println(maxSubstrings - result)
	return maxSubstrings - result
}
