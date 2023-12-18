package main

func main() {
	// fmt.Println(resp())
	solution("abacabad")
}

// func resp() int {
// 	array := []int{1, 2, 1, 3, 2, 5, 5, 6, 7, 8, 34, 8}
// 	m := make(map[int][]int)
// 	// given an array find the duplicates
// 	for i := 0; i < len(array); i++ {
// 		for j := i + 1; j < len(array); j++ {
// 			if array[i] == array[j] {
// 				m[array[i]] = append(m[array[i]], i)
// 				m[array[i]] = append(m[array[i]], j)
// 			}
// 		}
// 	}
// 	// compare the second index of each element to find the lower
// 	aux := 10000000000000000
// 	res := 0
// 	if len(m) == 0 {
// 		res = -1
// 		return res
// 	}
// 	for k, v := range m {
// 		fmt.Println(k, v)
// 		if v[1] < aux {
// 			aux = v[1]
// 			res = k
// 		}
// 	}
// 	return res
// }

/*
Given a string s consisting of small English letters, find and return the first instance of a non-repeating character in it. If there is no such character, return '_'.

Example

For s = "abacabad", the output should be
solution(s) = 'c'.

There are 2 non-repeating characters in the string: 'c' and 'd'. Return c since it appears in the string first.

For s = "abacabaabacaba", the output should be
solution(s) = '_'.

There are no characters in this string that do not repeat.
*/

func solution(s string) string {
	var arrayS []string

	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			if s[i] == s[j] {
				arrayS = append(arrayS, string(s[i]))
			}
		}
	}

	for i := 0; i < len(s); i++ {
		for j := 0; j < len(arrayS); j++ {
			if s[i] 
		}
	}
}
