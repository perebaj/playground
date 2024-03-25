package main

import "fmt"

func findPeaks(mountain []int) []int {
	var resp []int
	for i := 1; i < len(mountain)-1; i++ {
		if mountain[i] > mountain[i-1] && mountain[i] > mountain[i+1] {
			resp = append(resp, i)
		}
	}
	return resp
}

func main() {
	mountain := []int{2, 4, 4}
	fmt.Println(findPeaks(mountain))
}
