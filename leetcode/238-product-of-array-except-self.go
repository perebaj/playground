package main

import "fmt"

func main() {
	in := []int{1, 2, 3, 4}
	fmt.Println(productExceptSelf(in))
}

func productExceptSelf(nums []int) []int {
	sAux := make([]int, len(nums))

	for k := range sAux {
		aux := 1
		for kk, vv := range nums {
			if k != kk {
				aux = aux * vv
				sAux[k] = aux
			}
		}
	}
	fmt.Println(sAux)

	return sAux
}
