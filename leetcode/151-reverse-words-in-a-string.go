package main

import (
	"fmt"
	"strings"
)

func main() {
	s := "  hello world  "
	result := reverseWords(s)
	fmt.Println(result)
}

func reverseWords(s string) string {
	sliceS := strings.Fields(s)
	fmt.Println(sliceS)

	auxSSlice := make([]string, len(sliceS))
	for i := 0; i < len(sliceS); i++ {
		curIdx := (len(sliceS) - 1) - i
		auxSSlice[curIdx] = sliceS[i]
	}
	r := strings.Join(auxSSlice, " ")
	return r
}
