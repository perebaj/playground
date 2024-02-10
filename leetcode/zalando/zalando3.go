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


	   para resolver esse problema acho q podemos usar um hashmap para armazenar os dias que cada pessoa pode ir

	   e depois fazer um loop para ver quais dias tem mais pessoas disponiveis
	   o loop pode  ser

*/

package main

import (
	"fmt"
	"sort"
)

func main() {
	// fmt.Println(Solution([]string{"039", "4", "14", "32", "", "34", "7"})) // 5
	//
	fmt.Println(Solution([]string{"801234567", "180234567", "0", "189234567", "891234567", "98", "9"})) // 7
	// using 0 and 9
}

func Solution(E []string) int {
	// Implement your solution here
	m := make(map[int]int)
	for _, v := range E {
		for _, c := range v {
			m[int(c-'0')]++
		}
	}
	fmt.Println("map", m)
	var s []int
	for _, v := range m {
		s = append(s, v)
	}

	sort.Ints(s)
	fmt.Println("sorted array", s)
	var chaves []int
	for k, v := range m {
		if v == s[len(s)-1] || v == s[len(s)-2] {
			chaves = append(chaves, k)
		}
	}
	fmt.Println("chaves", chaves)
	var result int
	for _, v := range E {
		for _, c := range v {
			if int(c-'0') == chaves[0] || int(c-'0') == chaves[1] {
				result++
				break
			}
		}
	}

	return result
}
