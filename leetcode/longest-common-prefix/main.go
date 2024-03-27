package main

import (
	"fmt"
	"sort"
)

func main() {
	fmt.Println(longestCommonPrefix([]string{"flower", "flow", "flight"}))
}

func longestCommonPrefix(strs []string) string {
	sort.Strings(strs)
	fmt.Println(strs)
	var resp string
	for i := 0; i < len(strs[0]); i++ {
		fmt.Println(i)
		fmt.Println(strs[0][i], strs[len(strs)-1][i])
		if strs[0][i] == strs[len(strs)-1][i] {
			resp += string(strs[0][i])
		} else {
			return resp
		}

	}
	return resp
}
