package main

import "fmt"

func main() {
	word := "abcdefd"
	ch := byte('d')
	fmt.Println(reversePrefix(word, ch))
}

// It's possible to implement the logic to reverse and find the ch in the same loop
// But this approach that I'm using is O(N) time complexity
func reversePrefix(word string, ch byte) string {
	var toReverse string
	var index int
	for i := 0; i < len(word); i++ {
		if word[i] != ch {
			toReverse += string(word[i])
		} else {
			toReverse += string(word[i])
			index = i
			break
		}
	}
	fmt.Println("toReverse", toReverse)
	if len(toReverse) == len(word) && toReverse[index] != ch {
		return word
	}
	reversed := reverseString(toReverse)
	fmt.Println("reversedString", reversed)
	validStr := reversed + word[index+1:]
	return validStr

}

func reverseString(s string) string {
	var response string
	for i := len(s) - 1; i >= 0; i-- {
		response += string(s[i])
	}
	return response
}
