package main

import (
	"fmt"
)

func main() {
	fmt.Println(numJewelsInStones("z", "ZZ"))
}

func numJewelsInStones(jewels string, stones string) int {
	m := make(map[string]int)
	for i := 0; i < len(jewels); i++ {
		// fmt.Printf("%c\n", stones[i])
		m[string(jewels[i])] = 0
	}

	var result int

	for i := 0; i < len(stones); i++ {
		for k := range m {
			if string(stones[i]) == k {
				result++
			}
		}
	}

	return result
}

/*
   stones: is a representation about what I have
   jewels: is a representation of what is a valid gem

   - case sensitive

   1) Create a map where each key is a different jew
   2) compare with the stones and sum for each ocurrency
   	2.1 for loop over the stones
	2.2 verify if the current stones exists in the jew map
	2.3 sum in an aux var the got the result




*/
