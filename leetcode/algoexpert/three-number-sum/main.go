package main

import (
	"fmt"
	"sort"
)

func main() {
	array := []int{12, 3, 1, 2, -6, 5, -8, 6}
	targetSum := 0
	result := ThreeNumberSum(array, targetSum)
	fmt.Println(result)
}

// Não está correto, mas foi minha primeira ideia. Algumas solucões interessantes que encontrei aqui
func ThreeNumberSum(array []int, target int) [][]int{
	m1 := make(map[int]int)

	for _, v := range array {
		m1[v]++
	}
	fmt.Println(m1)
	m2 := make(map[int][]int)
	for i := 0; i < len(array); i++ {
		for j := i + 1; j < len(array); j++ {
			var sum int
			var auxArr []int
			sum = sum + array[i] + array[j]
			auxArr = append(auxArr, array[i])
			auxArr = append(auxArr, array[j])
			m2[sum] = auxArr
		}
	}

	for k, v := range m2 {
		_, ok := m1[(target - k)]
		if ok {
			fmt.Println(k, m1[(target-k)])
			m2[k] = append(v, (target - k))
		}
	}

	fmt.Println(m2)
	var results [][]int
	for _, v := range m2 {
		if len(v) == 3 {
			sort.Ints(v)
			results = append(results, v)
		}
	}

	results = removeDuplicates(results)
	return results
}

func removeDuplicates(arr [][]int) [][]int {
	unique := [][]int{}
	seen := map[string]bool{}

	for _, subArr := range arr {
		key := fmt.Sprint(subArr)
		if !seen[key] {
			unique = append(unique, subArr)
			seen[key] = true
		}
	}

	return unique
}
