/*
one-day-long training session will be conducted twice during the next 10 days.
There are N employees (numbered from 0 to N−1) willing to attend it. Each employee
has provided a list of which of the next 10 days they are able to participate in
the training. The employees’ preferences are represented as an array of strings. E[K]
is a string consisting of digits ('0'-'9') representing the days the K-th employee
is able to attend the training. The dates during which the training will take place
are yet to be scheduled. What is the maximum number of employees that can attend
during at least one of the two scheduled days? Write a function: func Solution(E []string)
int that, given an array E consisting of N strings denoting the available days for each
employee, will return the maximum number of employees that can attend during at least
one of the two scheduled days. Examples: 1. Given E = ["039", "4", "14", "32", "", "34", "7"],
the answer is 5. It can be achieved for example by running training on days 3 and 4.
This way employees number 0, 1, 2, 3 and 5 will attend the training. 2. Given
E = ["801234567", "180234567", "0", "189234567", "891234567", "98", "9"],
the answer is 7. It can be achieved by running training on days 0 and 9
. This allows all employees to attend the training. 3. Given E =
["5421", "245", "1452", "0345", "53", "354"], the answer is 6.
It can be achieved just by running training once on day 5, when every employee is available.

["039", "4", "14", "32", "", "34", "7"] // 5

each element in the array is trainer, and the content is the days they are available


0: 1
1: 1
2: 1
3: 3
4: 3
7: 1

{
	0: {
		0: 1
		3: 1
		9: 1
	},
	1: {
		4:1
	},
	2: {
		1: 1
		4: 1
	},
	3: {
		2: 1
		3: 1
	},
	4: {},
	5: {
		3: 1
		4: 1
	},
	6: {
		7: 1
	},
}
	// find what are the days that most are available
	// find overlap days in the data structure above

	most avaiable days: [3, 4] and the total is 6


}

*/

package main

import (
	"fmt"
)

func main() {
	fmt.Println(Solution([]string{"039", "4", "14", "32", "", "34", "7"})) // 5
	//
	// fmt.Println(Solution([]string{"801234567", "180234567", "0", "189234567", "891234567", "98", "9"})) // 7
	// using 0 and 9
}

func Solution(E []string) int {
	// Implement your solution here
	m := make(map[int]int)

	for _, v := range E {
		for _, vv := range v {
			string2int := (vv - '0')
			m[int(string2int)]++
		}
	}

	fmt.Println(m)
	var valuesSlice, keysSlice []int
	for k, v := range m {
		valuesSlice = append(valuesSlice, v)
		keysSlice = append(keysSlice, k)
	}
	fmt.Println(valuesSlice, keysSlice)

	max := valuesSlice[0]
	var keyMax int
	for k, v := range valuesSlice {
		if v > max {
			max = v
			keyMax = k
		}
	}
	fmt.Println(max, keyMax, keysSlice[keyMax])

	valuesSlice = remove(valuesSlice, keyMax)
	keysSlice = remove(keysSlice, keyMax)

	fmt.Println(valuesSlice, keysSlice)
	max2 := valuesSlice[0]
	var keyMax2 int
	for k, v := range valuesSlice {
		if v > max2 {
			max2 = v
			keyMax2 = k
		}
	}
	fmt.Println(max2, keyMax2, keysSlice[keyMax2])

	result := max2 + max

	for _, v := range v {
		
	}

	// fmt.Println()

	return 0
}

func remove(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
}
