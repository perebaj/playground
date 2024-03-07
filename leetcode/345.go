package main

import "fmt"

func main() {
	s := "leetcode"
	// leetcedo
	fmt.Println(reverseVowels(s))
}

func reverseVowels(s string) string {
	s2 := []rune(s)
	var vowelCount int
	for _, value := range s {
		switch value {
		case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
			vowelCount++
		}
	}
	fmt.Println(vowelCount)
	var index int
	var v rune
	for {
		if isVowel(s2[index]) && v == 0 {
			v = s2[index]
		} else if isVowel(s2[index]) {
			aux := s2[index]
			s2[index] = v
			v = aux
			vowelCount--
		}

		if index == len(s)-1 {
			index = 0
		} else {
			index++
		}
		if vowelCount == 0 {
			break
		}
	}
	fmt.Println(string(s2))
	return string(s2)
}

func isVowel(c rune) bool {
	switch c {
	case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
		return true
	default:
		return false
	}
}

// can appear in both lower and upper cases,
/*

Example 1:

Input: s = "hello"
Output: "holle"
Example 2:

Input: s = "leetcode"
Output: "leotcede"


1) identify how many vowels
2) iterate over the string decreasing the vowels var until it achieve 0


nVowlds := 4

for {
	if VOLWEL && v == ""{
		v = string[index]
	} else {
		aux := string[index]
		string[index] = v
		v = aux
		nVolwels--
	}

	if index == len(s) -1 {
		index = 0
	}

	if nVolwels == 0 {
		break
	}
}

*/
