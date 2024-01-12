package main

import "fmt"

func main() {
	in := [][]string{
		{"qMTSlfgZlC", "ePvzZaqLXj"},
		{"xKhZXfuBeC", "TtnllZpKKg"},
		{"ePvzZaqLXj", "sxrvXFcqgG"},
		{"sxrvXFcqgG", "xKhZXfuBeC"},
		{"TtnllZpKKg", "OAxMijOZgW"},
	}
	fmt.Println(in[0])
	fmt.Println(destCity(in))
}

func destCity(paths [][]string) string {
	m1, m2 := make(map[string]int), make(map[string]int)

	for i := 0; i < len(paths); i++ {
		m1[paths[i][0]]++
		m2[paths[i][1]]++
	}

	var resp string
	for k := range m2 {
		_, ok := m1[k]
		if !ok {
			resp = k
		}
	}

	return resp
}
