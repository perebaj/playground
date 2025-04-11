package main

import "fmt"

func main() {
	m := make(map[string]int)

	m["key1"] = 1

	v, ok := m["key2"]

	fmt.Println("ok", ok)
	fmt.Println("v", v)
}
