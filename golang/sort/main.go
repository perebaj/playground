package main

import (
	"fmt"
	"sort"
)

type Person struct {
	Name string
	Age  int
}

type People []Person

func (p People) Len() int {
	return len(p)
}

func (p People) Less(i, j int) bool {
	if p[i].Age == p[j].Age { // if the age is the same, sort by name, otherwise, sort by age
		// where i is the current index and j is the next index
		return p[i].Name < p[j].Name
	}

	return p[i].Age > p[j].Age // if we want to sort using the ASC or DESC, we just need to change the > operator to < and vice-versa
}

// Swap is used to swap the elements in the slice
func (p People) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func main() {
	people := People{
		{"Alice", 25},
		{"Bob", 30},
		{"Marlon", 25},
		{"David", 20},
	}

	// we can user the sort.Sort passing a slice of the struct, because we implement our version of the sort.Interface
	// Preatty cool, right?
	sort.Sort(people)

	fmt.Println(people)
}
