package main

import "fmt"

// Here we are passing a copy of the array, for this reason the original value won't be modified, and the value of the array will be the same
// remembering if you don't determine a value, the zero-value will be assigned, for a slice of ints, will be 0
func update(x [3]int) {
	x[2] = 5
}

func updateString(x [3]string) {
	x[2] = "5"
}

func main() {
	var a = [3]int{1, 2}
	update(a)
	fmt.Println(a)
	var b = [3]string{"1", "2"}
	updateString(b)
	fmt.Println(b)
}
