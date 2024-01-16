package main

import (
	"fmt"
	"math"
	"strconv"
)

func main() {
	number := 2
	fmt.Println(isHappy(number))
}

func isHappy(n int) bool {
	numberString := strconv.Itoa(n)
	m := make(map[float64]int)
	var aux float64
	for {
		for _, v := range numberString {
			i, _ := strconv.Atoi(string(v))
			aux = aux + math.Pow(float64(i), 2)
		}

		if aux == 1 {
			break
		} else {
			_, ok := m[aux]
			if ok {
				return false
			}
			m[aux] = 0
			numberString = strconv.Itoa(int(aux))
			aux = 0

		}
	}
	return true
}
