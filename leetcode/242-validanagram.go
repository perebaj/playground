package main

import "fmt"

func main() {
	s := "aacc"
	t := "ccac"
	fmt.Println(isAnagram(s, t))
}

func isAnagram(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}

	sMap := make(map[byte]int)
	for i := 0; i < len(s); i++ {
		sMap[s[i]]++
	}

	tMap := make(map[byte]int)
	for i := 0; i < len(t); i++ {
		tMap[t[i]]++
	}

	for i := 0; i < len(s); i++ {
		v, ok := tMap[s[i]]
		if !ok || sMap[s[i]] != v {
			return false
		}
	}

	return true
}

/*
	create an hashmap using the current statement

	1) len of both are the same
	2) create an hashmap using some of the word letters as keys in HM
	3) count the occurencies of letter in the second word in the above HM created
	in: eat
	{
		e:
		a:
		t:
	}

	for over t string{
		_, ok := m[t]
		if !ok {
			return false
		} else {
			m[t]++
		}
	}

*/
