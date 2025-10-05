package main

import (
	"fmt"

	"gonum.org/v1/gonum/stat"
)

func main() {
	arr1 := []float64{19, 12, 13, 0}
	arr2 := []float64{100, 28, 71, 6}
	fmt.Println(stat.LinearRegression(arr1, arr2, nil, false))

}
