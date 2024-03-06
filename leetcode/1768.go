package main

import "fmt"

func main() {
	word1 := "abc"
	word2 := "pq"
	fmt.Println(mergeAlternately(word1, word2))
}

func mergeAlternately(word1 string, word2 string) string {
	len1 := len(word1)
	len2 := len(word2)
	var i int
	var result string
	for {
		if i < len1 {
			result += string(word1[i])
		}
		if i < len2 {
			result += string(word2[i])
		}
		if i > len1 && i > len2 {
			break
		}
		fmt.Println(string(result))
		i++
	}
	return result
}
