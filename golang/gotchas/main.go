package main

import (
	"fmt"
	"math"
	"time"
)

type Point struct {
	X, Y float64
}

func (p *Point) Abs() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

func main() {
	var p *Point
	// 1

	fmt.Println(p) // it's possible to verify here, that the value of p is nil, this because when you create a pointer, the default value
	// assigned to it, it's nil. For this reason, when we try to access the X and Y, in the Abs function, isn't possible

	p = new(Point) // In this case we need to intialize the variabel
	fmt.Println(p.X)

	// 2
	// the same occour with maps

	var m map[int]string
	// m[1] = "jonathan" <- this line will generate an error

	m = make(map[int]string)
	m[1] = "jonathan"
	fmt.Println(m)

	// 3
	// t := time.Parse(time.RFC3339, "2018-04-06T10:49:05Z") //assignment mismatch: 1 variable but time.Parse returns 2 values

	// The above line will return an error, because the functiosn returns 2 variables, the right way to code it, is:

	t, _ := time.Parse(time.RFC3339, "2018-04-06T10:49:05Z")
	fmt.Println(t)

	// 4

	a := [2]int{1, 2} // array definition
	changeArray(a)
	fmt.Println(a)
	b := []int{1, 2}
	changeSlice(b)
	fmt.Println(b)
}

func changeArray(a [2]int) {
	a[1] = 33
}

func changeSlice(a []int) {
	a[1] = 55
}
