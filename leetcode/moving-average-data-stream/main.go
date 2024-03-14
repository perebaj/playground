package main

import "fmt"

type MovingAvg struct {
	Avg   float32
	Count int
}

func main() {
	m := &MovingAvg{}
	m.Next(10.5)
	m.Next(2)

	fmt.Println(m.Avg) // 6
}

func (m *MovingAvg) Next(in float32) {
	m.Count++
	m.Avg = (m.Avg + in) / float32(m.Count)
}
