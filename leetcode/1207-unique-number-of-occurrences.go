package main

import "fmt"

func main() {
	arr := []int{1, 2, 2, 1, 1, 3}
	fmt.Println(uniqueOccurrences(arr))

	// m := make(map[int]int)

	// m[1] = 0

	// _, ok := m[2]
	// fmt.Println(ok)
}

func uniqueOccurrences(arr []int) bool {
	m := make(map[int]int)

	for i := 0; i < len(arr); i++ {
		m[arr[i]]++
	}

	// verify if the hm values are unique
	//[3, 2, 1]
	var auxSlice []int
	for _, v := range m {
		auxSlice = append(auxSlice, v)
	}

	auxM := make(map[int]int)
	for _, v := range auxSlice {
		_, ok := auxM[v]
		if !ok { //false
			auxM[v] = 0
		} else {
			return false
		}
	}
	return true
}
