package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

/*
	1) recognize all the number in the string
		string = len(string) = S = O(S)
	2) find the first and the last
		2.1) if numbers == 1. Then, this number is the fast and the last
		O(1)
	3) Do that on all lines
		all lines = size of F = O(F)
	4) sum all ocurencies and return the result


	O(S) * O(F) = O(nˆ2)
*/

func main() {
	f, err := os.Open("./1-input.txt")
	if err != nil {
		panic("failed to read")
	}

	r := bufio.NewReader(f)
	var sum int
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}

		result := findNumber(line)
		var calibrationVal string
		if len(result) < 2 {
			calibrationVal = string(result[0]) + string(result[0])
		}
		calibrationVal = string(result[0]) + string(result[len(result)-1])

		calibrationValueInt, _ := strconv.Atoi(calibrationVal)
		sum = sum + calibrationValueInt
	}
	fmt.Println(sum)
}

func findNumber(s string) []rune {
	var result []rune
	for _, v := range s {
		ok := unicode.IsNumber(v)
		if ok {
			result = append(result, v)
		}
	}
	return result
}
