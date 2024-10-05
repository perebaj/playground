package main

import "fmt"

func main() {
	array := []int{2, 1, 2, 2, 2, 3, 4, 2}
	toMove := 2

	res := MoveElementToEnd(array, toMove)
	fmt.Println(res)
}

func MoveElementToEnd(array []int, toMove int) []int {
	// Write your code here.
	l, r := 0, 1

	for r < len(array) {
		// swap rule
		if array[l] == toMove && array[r] != toMove {
			aux := array[l]
			array[l] = array[r]
			array[r] = aux
			l++
			r++
		} else if array[r] == toMove && array[l] != toMove {
			l = r
			r++
		} else {
			r++
		}
	}

	return array
}
