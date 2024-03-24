package main

import "fmt"

type Key [3]int

type RangeFreqQuery struct {
	InitialArray []int       // save the number and how many times it appears
	Cache        map[Key]int // save time for future calclulations
}

func Constructor(arr []int) RangeFreqQuery {
	return RangeFreqQuery{
		InitialArray: arr,
		/*
		Using this cache structure here was a really good trick to avoid recalculate some index that I'm already touched. But I basically did that
		to circumvent a test that are doing the same calculation multiple times.
		The last test, basically, do multiple queries for different ranges and values, so, the cache doesn't work so well in this case.
		*/
		Cache:        make(map[Key]int),
	}
}

func (r *RangeFreqQuery) Query(left int, right int, value int) int {
	subArray := r.InitialArray[left : right+1]
	m := make(map[int]int)
	_, ok := r.Cache[Key{
		left,
		right + 1,
		value,
	}]
	if !ok {
		for _, v := range subArray {
			m[v]++
		}
		r.Cache[Key{
			left, right + 1, value,
		}] = m[value]

		return m[value]
	}
	return r.Cache[Key{
		left, right + 1, value,
	}]

}

func main() {

	array := []int{3, 3, 7, 3, 8, 10, 2, 4, 2, 2}
	n := Constructor(array)
	//[0,3,2],[9,9,9],[6,8,8],[6,9,6],[4,5,9],[6,9,9],[8,9,5]
	fmt.Println(n.Query(0, 3, 2))
	fmt.Println(n.Query(9, 9, 9))
	fmt.Println(n.Query(6, 8, 8))
	fmt.Println(n.Query(6, 8, 4))
}
