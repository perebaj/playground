package main

import (
	"fmt"
	"reflect"
)

func main() {
	mainValue := []int{49, 18, 16}
	sub := [][]int{
		{16, 18, 49},
	}
	fmt.Println(solution(mainValue, sub))

	mainValue = []int{15, 88, 88}
	sub = [][]int{
		{88},
		{15},
	}
	fmt.Println(solution(mainValue, sub))
}

func solution(mainValue []int, sub [][]int) bool {
	m := make(map[int]int)
	for _, v := range mainValue {
		m[v]++
	}

	for k, v := range sub {
		if len(v) > 1 && !reflect.DeepEqual(mainValue[k:len(sub)], sub) {
			fmt.Println("primeiro caso")
			return false
		} else {
			fmt.Println("segundo caso")
			_, ok := m[v[0]]
			if !ok {
				fmt.Println("segundo.1 caso")
				return false
			} else {
				fmt.Println("segundo.2 caso")
				continue
			}
		}
	}

	return true
}
