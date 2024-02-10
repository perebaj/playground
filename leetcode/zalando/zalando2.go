/*
A string S consisting of the letters A, B, C and D is given. The string can be
transformed either by removing a letter A together with an adjacent letter B, or
 by removing a letter C together with an adjacent letter D. Write a function:
 func Solution(S string) string that, given a string S consisting of N characters,
 returns any string that: can be obtained from S by repeatedly applying the described
  transformation, and cannot be further transformed. If at some point there is more than
  one possible way to transform the string, any of the valid transformations may be chosen.

Examples:
1. Given S = "CBACD", the function may return "C"

BA -> CCD -> C

2) Given S = "CABABD", the function may return "C"

AB -> CABD -> CD -> C

3) Given S = "ACBDACBD", the function may return "ACBD"


CREATE MULTIPLE TEST CASES

ABCDACBD
BBBBBCC


*/

package main

import "fmt"

func main() {

	// fmt.Println(solution("CBACD"))

	// fmt.Println(solution("CABABD"))

	fmt.Println(solution("B"))

}

func delChar(s []rune, index int) []rune {
	return append(s[0:index], s[index+1:]...)
}

func solution(S string) string {
	if len(S) == 0 {
		return ""
	}
	if len(S) < 2 {
		return S
	}
	var index int
	for {
		// fmt.Println("S: ", S)
		// fmt.Println("index: ", index)
		// fmt.Println("TESTE: ", S[index:2+index])
		if S[index:2+index] == "AB" || S[index:2+index] == "CD" || S[index:2+index] == "BA" || S[index:2+index] == "DC" {
			// fmt.Println("ENTROU")
			S = S[0:index] + S[index+2:]
			index = 0
			if len(S) < 2 {
				break
			}
		} else {
			index++
			if index == len(S)-1 {
				break
			}
		}
	}

	return S
}
