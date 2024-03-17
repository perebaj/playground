package main

import "fmt"

func main() {
	s := "axc"
	t := "ahbgdc"
	fmt.Println(isSubsequence(s, t))
}

func isSubsequence(s string, t string) bool {
	if len(s) < 1 {
		return true
	}
	position := 0
	var result int
	for i := 0; i < len(t); i++ {
		if s[position] == t[i] {
			result++
			position++
		}
		if position == len(s) {
			break
		}
	}

	if result != len(s) {
		return false
	}
	return true
}
