package main

import "fmt"

func main() {
	fmt.Println(reverseVowels("hello"))
}

func reverseVowels(s string) string {
	var reversedVowels string
	var reversedVowelsIndex []int
	for i, v := range s {
		if isVowel(v) {
			reversedVowelsIndex = append(reversedVowelsIndex, i)
			reversedVowels = string(v) + reversedVowels
		}
	}
	if len(reversedVowels) < 2 {
		return s
	}

	sByte := []byte(s)
	var pointer int
	for i, _ := range sByte {
		if pointer < len(reversedVowelsIndex) && i == reversedVowelsIndex[pointer] {
			sByte[i] = reversedVowels[pointer]
			pointer++
		}
	}

	fmt.Println("reversedVowels", reversedVowels)
	fmt.Println("sByte", sByte)

	return string(sByte)
}

func isVowel(c rune) bool {
	switch c {
	case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
		return true
	default:
		return false
	}
}
