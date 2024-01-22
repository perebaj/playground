package main

import "fmt"

func main() {
	slice := []int{1, 2, 3, 4, 5}

	a := fmt.Sprintf("%v", slice)
	fmt.Printf("%T, %v\n", a, a)

}
