package main

import (
	"fmt"
	"sort"
)

/*
The right thing to do is to start choosing the jobs by the highests deadline from the lowest.

For example. Jobs with deadline 3 can compete with the slot 1,2 and 3. So, we need to compare them all.

Example

d 1 p 1
d 2 p 2
d 3 p 3
d 3 p 4
d 3 p 5

the best chooses will be all payments with deedline equal to 3
*/

func main() {
	jobs := []map[string]int{
		//case 1
		{"deadline": 1, "payment": 2},
		{"deadline": 1, "payment": 1},
		{"deadline": 1, "payment": 3},
		{"deadline": 2, "payment": 2},
		{"deadline": 2, "payment": 3},
		{"deadline": 3, "payment": 3},
		{"deadline": 4, "payment": 4},
		{"deadline": 4, "payment": 2},
		{"deadline": 5, "payment": 5},

		//case 2
		// {"deadline": 2, "payment": 1},
		// {"deadline": 2, "payment": 2},
		// {"deadline": 2, "payment": 3},
		// {"deadline": 2, "payment": 4},
		// {"deadline": 2, "payment": 5},
		// {"deadline": 2, "payment": 6},
		// {"deadline": 2, "payment": 7},
	}

	// sort.Slice(jobs, func(i, j int) bool {
	// 	fmt.Println(jobs[i]["deadline"], jobs[j]["deadline"])
	// 	if jobs[i]["deadline"] == jobs[j]["dealine"] {
	// 		return jobs[i]["payment"] > jobs[j]["payment"]
	// 	}
	// 	return jobs[i]["deadline"] < jobs[j]["dealine"]
	// })

	sort.Slice(jobs, func(i, j int) bool {
		if jobs[i]["payment"] == jobs[j]["payment"] {
			return jobs[i]["deadline"] > jobs[j]["deadline"]
		}
		return jobs[i]["payment"] > jobs[j]["payment"]
	})

	fmt.Println(jobs)
}
