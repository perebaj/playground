// https://leetcode.com/discuss/interview-question/406652/Twitter-or-OA-2019-or-Anagram-Difference
package main

import "fmt"

func main() {
	fmt.Println("asdklajs")
	a := []string{"tea", "tea", "a", "jk", "abb", "mn", "abc"}
	b := []string{"ate", "eat", "bb", "kj", "bbc", "op", "def"}
	// getMinimumDifference(a, b)
	fmt.Println(getMinimumDifference(a, b))
}

func individualDif(a string, b string) int {
	if len(a) != len(b) {
		return -1
	}

	mapA := make(map[rune]int)
	mapB := make(map[rune]int)

	for _, v := range a {
		mapA[v]++
	}

	for _, v := range b {
		mapB[v]++
	}

	fmt.Println(mapA, mapB)

	var resp int
	for k, v := range mapA {
		if mapB[k] == v {
			continue
		} else {
			resp += v
		}
	}
	return resp
}

func getMinimumDifference(a []string, b []string) []int32 {
	var fResp []int32

	for i := 0; i < len(a); i++ {
		resp := individualDif(a[i], b[i])
		fResp = append(fResp, int32(resp))
	}
	return fResp
}
