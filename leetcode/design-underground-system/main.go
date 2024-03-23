package main

import "fmt"

type UndergroundSystem struct {
	Cities      map[string][]int
	CheckInMap  map[int]Check
	CheckOutMap map[int]Check
}

type Check struct {
	Name string
	Time int
}

func Constructor() UndergroundSystem {
	return UndergroundSystem{
		Cities:      make(map[string][]int),
		CheckInMap:  make(map[int]Check),
		CheckOutMap: make(map[int]Check),
	}
}

func (u *UndergroundSystem) CheckIn(id int, stationName string, t int) {
	c := Check{
		Name: stationName,
		Time: t,
	}
	u.Cities[stationName] = append(u.Cities[stationName], id)
	u.CheckInMap[id] = c
}

func (u *UndergroundSystem) CheckOut(id int, stationName string, t int) {
	c := Check{
		Name: stationName,
		Time: t,
	}

	u.CheckOutMap[id] = c
}

func (u *UndergroundSystem) GetAverageTime(startStation string, endStation string) float64 {
	ids := u.Cities[startStation]

	var sliceTimes []int

	for _, v := range ids {
		in := u.CheckInMap[v]
		out := u.CheckOutMap[v]
		if out.Name == endStation {
			time := out.Time - in.Time
			sliceTimes = append(sliceTimes, time)
		} else {
			continue
		}
	}
	var sum int
	for _, v := range sliceTimes {
		sum += v
	}
	result := float64(sum) / float64(len(sliceTimes))
	return result
}

func main() {
	n := Constructor()
	n.CheckIn(10, "Leyton", 3)
	n.CheckOut(10, "Paradise", 8)
	fmt.Println(n.GetAverageTime("Leyton", "Paradise")) // return 5.00000, (5) / 1 = 5
	n.CheckIn(5, "Leyton", 10)
	n.CheckOut(5, "Paradise", 16)
	fmt.Println(n.CheckInMap, n.CheckOutMap)
	fmt.Println(n.GetAverageTime("Leyton", "Paradise")) // return 5.50000, (5 + 6) / 2 = 5.5

	// n.CheckIn(1, "ocz", 10)
	// n.CheckOut(1, "parapua", 20)
	// n.GetAverageTime("ocz", "parapua")
}
