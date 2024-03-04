package main

import (
	"fmt"
	"strings"
)

func main() {
	pattern := "abba"
	s := "dog cat cat dog"
	fmt.Println(wordPattern(pattern, s))
}

func wordPattern(pattern, s string) bool {
	sliceS := strings.Split(s, " ") // slice of S
	fmt.Println(sliceS)
	if len(sliceS) != len(pattern) {
		return false
	}
	patternHM := make(map[string]string)
	SHM := make(map[string]string)

	for i := 0; i < len(sliceS); i++ {
		vpatternHM, ok := patternHM[sliceS[i]]
		vSHM, ok2 := SHM[string(pattern[i])]

		if !ok && !ok2 {
			patternHM[sliceS[i]] = string(pattern[i])
			SHM[string(pattern[i])] = sliceS[i]
		} else if vpatternHM != SHM[string(pattern[i])] && vSHM != patternHM[sliceS[i]] {
			return false
		}
		fmt.Println(vpatternHM, vSHM)
	}
	fmt.Println(patternHM, SHM)
	return true
}

// func wordPattern(pattern, s string) bool {
// 	sliceS := strings.Split(s, " ")
// 	fmt.Println(sliceS)
// 	if len(sliceS) != len(pattern) {
// 		return false
// 	}
// 	m := make(map[string]string)

// 	for i := 0; i < len(sliceS); i++ {
// 		vPattern, ok := m[string(pattern[i])]
// 		vSlice, ok2 := m[sliceS[i]]

// 		if !ok && !ok2 {
// 			m[string(pattern[i])] = sliceS[i]
// 			m[sliceS[i]] = string(pattern[i])
// 		} else if vPattern != sliceS[i] && vSlice != string(pattern[i]) {
// 			return false
// 		}
// 	}
// 	fmt.Println(m)
// 	return true
// }

// func wordPattern(pattern string, s string) bool {
// 	//split

// 	sliceS := strings.Split(s, " ")
// 	fmt.Println(sliceS)
// 	if len(sliceS) != len(pattern) {
// 		return false
// 	}

// 	m := make(map[string]string)

// 	pKey := string(pattern[0])
// 	pValue := sliceS[0]
// 	m[pKey] = pValue
// 	for i := 1; i < len(pattern); i++ {
// 		vv, ok := m[string(pattern[i])]
// 		fmt.Println(ok, sliceS[i], pValue)
// 		if !ok && sliceS[i] != pValue {
// 			m[string(pattern[i])] = sliceS[i]
// 			pValue = sliceS[i]
// 		} else if vv != sliceS[i] {
// 			return false
// 		}
// 	}

// 	return true
// }
