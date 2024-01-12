package main

import (
	"fmt"
	"sort"
)

func main() {
	s := "aaaabbbbcccc"
	fmt.Println(sortString(s))
}

func sortString(s string) string {
	m := make(map[int32]int)

	for _, v := range s {
		m[v]++
	}

	// 2) coloco os valores das chaves em um array sorted?
	//     obs: para conseguir seguir as regras de 1-6
	var auxSlice []int32
	for k := range m {
		auxSlice = append(auxSlice, k)
	}

	sort.Slice(auxSlice, func(i, j int) bool { return auxSlice[i] < auxSlice[j] })

	var resp string
	for {
		for j := 0; j < len(auxSlice); j++ {
			if m[auxSlice[j]] > 0 {
				resp = resp + string(auxSlice[j])
				m[auxSlice[j]]--
			}
		}

		for i := len(auxSlice) - 1; i >= 0; i-- {
			if m[auxSlice[i]] > 0 {
				resp = resp + string(auxSlice[i])
				m[auxSlice[i]]--
			}
		}
		var sum int
		for _, v := range m {
			sum = sum + v
		}
		if sum == 0 {
			break
		}
	}

	return resp
}
