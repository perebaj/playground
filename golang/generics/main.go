/*
Tutorial: https://go.dev/doc/tutorial/generics
Cheat sheet: https://gosamples.dev/generics-cheatsheet/
*/
package main

import "fmt"

func main() {
	ints := map[string]int64{"first": 1, "second": 2, "third": 3}

	floats := map[string]float64{"first": 1.1, "second": 2.2, "third": 3.3}

	fmt.Printf("non-generic sums: %v and %v\n", SumInsts(ints), SumFloats(floats))

	fmt.Printf("generic sum: %v\n", SumIntsOrFloats[string, int64](ints))
	fmt.Printf("generic sum: %v\n", SumIntsOrFloats[string, float64](floats))

	fmt.Printf("generic sum: %v\n", SumNumbers(ints))
	fmt.Printf("generic sum: %v\n", SumNumbers(floats))
}

func SumInsts(m map[string]int64) int64 {
	var sum int64
	for _, v := range m {
		sum += v
	}
	return sum
}

func SumFloats(m map[string]float64) float64 {
	var sum float64
	for _, v := range m {
		sum += v
	}
	return sum
}

type Number interface {
	int64 | float64
}

func SumIntsOrFloats[k comparable, V int64 | float64](m map[string]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

func SumNumbers[K comparable, V Number](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}
