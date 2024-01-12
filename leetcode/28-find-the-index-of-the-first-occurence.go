package main

import "fmt"

func main() {
	haystack := "hello"
	needle := "ll"

	fmt.Println(strStr(haystack, needle))
}

func strStr(haystack string, needle string) int {
	needleLen := len(needle)
	for i := 0; i < len(haystack)-len(needle)+1; i++ {
		if haystack[i:i+needleLen] == needle {
			return i
		}
	}
	return -1
}
