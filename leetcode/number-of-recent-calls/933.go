package main

import "fmt"

type RecentCounter struct {
	requests []int
}

func Constructor() RecentCounter {
	return RecentCounter{}
}

func (this *RecentCounter) Ping(t int) int {
	this.requests = append(this.requests, t)
	rangeSlice := Range(t)
	var counter int
	for _, v := range this.requests {
		if v >= rangeSlice[0] && v <= rangeSlice[1] {
			counter++
		}
	}
	return counter
}

func Range(t int) []int {
	range1 := t - 3000
	range2 := t
	resp := make([]int, 2)
	if range1 > range2 {
		resp[0] = range2
		resp[1] = range1
	} else {
		resp[0] = range1
		resp[1] = range2
	}

	return resp
}

func main() {
	obj := Constructor()
	fmt.Println(obj.Ping(1))
	fmt.Println(obj.Ping(100))
	fmt.Println(obj.Ping(3001))
	fmt.Println(obj.Ping(3002))
}
