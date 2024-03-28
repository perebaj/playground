package main

func main() {
	words := []string{"abc", "xyx", "aba", "1221"}
	println(firstPalindrome(words))
}

func firstPalindrome(words []string) string {
	for _, value := range words {
		r := isPalindrome(value)
		if r {
			return value
		}
	}
	return ""
}

func isPalindrome(word string) bool {
	start := 0
	end := len(word) - 1
	for start < end {
		if word[start] != word[end] {
			return false
		}
		start++
		end--
	}
	return true
}
