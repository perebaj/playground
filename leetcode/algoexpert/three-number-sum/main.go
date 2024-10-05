package main

import (
	"fmt"
	"sort"
)

func main() {
	array := []int{1, 2, 3}
	targetSum := 6
	result := ThreeNumberSum2(array, targetSum)
	fmt.Println(result)
}

func ThreeNumberSum2(array []int, target int) [][]int {
	sort.Ints(array)
	var result [][]int
	if len(array) < 3 {
		return result
	}
	end := len(array) - 1
	for i := 0; i < len(array)-3; i++ {
		green := i + 1
		blue := end - 1
		for {
			sum := array[i] + array[green] + array[blue]
			fmt.Println(sum)
			if sum == target {
				result = append(result, []int{array[i], array[green], array[blue]})
			}
			blue--
			if green == end {
				break
			}

			if green == blue {
				green++
				blue = end
			}
		}
	}
	return result
}

// Não está correto, mas foi minha primeira ideia. Algumas solucões interessantes que encontrei aqui
func ThreeNumberSum(array []int, target int) [][]int {
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
