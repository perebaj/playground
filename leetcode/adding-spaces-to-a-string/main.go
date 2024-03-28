package main

import "fmt"

func main() {
	s := "HelloWorld"
	spaces := []int{5}
	fmt.Println(addSpaces(s, spaces))

}

// I don't know why this solution exceeds the time limit, I consider it O(N) time complexity
func addSpaces(s string, spaces []int) string {
	for key, value := range spaces {
		s = s[0:key+value] + " " + s[key+value:]
		// fmt.Println(s)
	}
	return s
}
