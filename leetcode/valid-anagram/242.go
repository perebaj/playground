package main

import "fmt"

func main() {
	s := "anagrmm"
	t := "nagaram"
	fmt.Println(isAnagram(s, t))
}

func isAnagram(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}

	mapS := make(map[byte]int)
	mapT := make(map[byte]int)

	for i := 0; i < len(s); i++ {
		mapS[s[i]]++
		mapT[t[i]]++
	}

	for kS, vS := range mapS {
		vT, ok := mapT[kS]
		if !ok || vS != vT {
			return false
		}
	}

	return true
}
