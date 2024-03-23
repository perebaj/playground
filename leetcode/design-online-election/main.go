package main

import (
	"fmt"
	"sort"
)

type TopVotedCandidate struct {
	timePersonIdMap map[int]int
	sortedTimes     []int
}

/*
Complexity Analysis

first for: O(N) where N is the size of persons array
second for to store all the times in a slice: O(N)
sort the times: O(NlogN)

I.e: O(N) + O(N) + O(NlogN), approximately O(NlogN)
*/
func Constructor(persons []int, times []int) TopVotedCandidate {
	maxVotes := -1
	m := make(map[int]int)  // person -> votes frequency ex: 1:2, 2:3. Id 1 has 2 votes, id 2 has 3 votes
	m2 := make(map[int]int) // time -> person id ex: 0:1, 5:1 at time 0 the leader is person 1, at time 5 the leader is person 1
	// var leader int
	for i := 0; i < len(persons); i++ {
		m[persons[i]]++
		if m[persons[i]] >= maxVotes {
			maxVotes = m[persons[i]]
			m2[times[i]] = persons[i]
		}
	}

	var sortedTimes []int
	for k := range m2 {
		sortedTimes = append(sortedTimes, k)
	}

	// sort the times
	sort.Ints(sortedTimes)

	return TopVotedCandidate{
		timePersonIdMap: m2,
		sortedTimes:     sortedTimes,
	}

}

/*
	Time to search the target element in a sorted array: O(logN)
*/
func (top *TopVotedCandidate) Q(t int) int {
	// var result []int

	val := sort.Search(len(top.sortedTimes), func(i int) bool {
		return top.sortedTimes[i] >= t
	})
	// if the val is in this range, this means that the target exists in the sortedTimes, otherwise, the function return the index where the target should be inserted
	// in this specific case, the previous index.
	if val < len(top.sortedTimes) && top.sortedTimes[val] == t {
		return top.timePersonIdMap[t]
	}
	return top.timePersonIdMap[top.sortedTimes[val-1]]
}


/*
Final complexity analysis
the constructor has O(NlogN) complexity
the query has O(logN) complexity

I.e: O(NlogN) + O(logN) = O(NlogN)
*/
func main() {
	person := []int{0, 1, 1, 0, 0, 1, 0}
	times := []int{0, 5, 10, 15, 20, 25, 30}
	obj := Constructor(person, times)
	fmt.Println(obj.timePersonIdMap)
	fmt.Println(obj.Q(3))
	fmt.Println(obj.Q(12))
	fmt.Println(obj.Q(25))
	fmt.Println(obj.Q(15))
	fmt.Println(obj.Q(24))
	fmt.Println(obj.Q(8))
}
