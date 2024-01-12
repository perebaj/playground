package main

import (
	"fmt"
	"strings"
)

func main() {
	s := "   fly me   to   the moon  "
	fmt.Println(lengthOfLastWord(s))
}

func lengthOfLastWord(s string) int {
	s = strings.TrimSpace(s)
	result := strings.Split(s, " ")
	return len(result[len(result)-1])
}
