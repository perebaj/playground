package main

import "fmt"

func lengthOfLongestSubstring(s string) int {
	m := make(map[rune]int)
	counter := 0
	var auxCounter int
	for i := 0; i < len(s); i++ {
		for j := i; j < len(s); j++ {
			_, ok := m[rune(s[j])]
			// not exist
			if !ok {
				m[rune(s[j])]++
				auxCounter++
				if auxCounter > counter {
					counter = auxCounter
				}
			} else {
				// reset the map and set i to be equal to j
				m = make(map[rune]int)
				auxCounter = 0
				break
			}
		}
	}

	return counter
}

func main() {
	fmt.Println(lengthOfLongestSubstring("pwwkew"))
}
