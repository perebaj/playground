package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	Day2Part1()
	Day2Part2()
}

func Day2Part1() {
	elements := readInput("./day2.txt")
	// fmt.Println(elements)

	var results int
	for _, v := range elements {
		if isMonotonic(v) {
			results++
		}
	}
	fmt.Println("Day 2 Part 1 result: ", results)
}

func Day2Part2() {
	elements := readInput("./day2.txt")

	var results int
	for _, v := range elements {
		if !isMonotonic(v) {
			for i := 0; i < len(v); i++ {
				tmp := make([]string, len(v))
				copy(tmp, v)
				aux := slices.Delete(tmp, i, i+1)
				if isMonotonic(aux) {
					results++
					break
				}
			}
		} else {
			results++
		}
	}

	fmt.Println("Day 2 Part 2 result: ", results)
}

func readInput(path string) [][]string {
	file, err := os.Open(path)
	if err != nil {
		panic("failed to read")
	}

	defer file.Close()

	var elements [][]string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		elementsLine := strings.Split(line, " ")

		elements = append(elements, elementsLine)
	}

	if err := scanner.Err(); err != nil {
		fmt.Errorf("%v", err)
	}

	return elements
}

// isMonotonic validates if some rules given an slice of string
// The elements are either all increasing or all decreasing.
// Any two adjacent elements differ by at least one and at most three.
func isMonotonic(elements []string) bool {
	decreasing := false
	increasing := false

	// convert []string to []int
	var elementsInt []int
	for _, v := range elements {
		i, _ := strconv.Atoi(v)
		elementsInt = append(elementsInt, i)
	}

	for i := 0; i < len(elementsInt)-1; i++ {
		if elementsInt[i] > elementsInt[i+1] {
			increasing = true
		} else if elementsInt[i] < elementsInt[i+1] {
			decreasing = true
		} else if elementsInt[i] == elementsInt[i+1] {
			return false
		}

		// if the diference between element and element+1 was bigger than 3
		// return false
		absDiff := absDif(elementsInt[i], elementsInt[i+1])
		if absDiff > 3 {
			return false
		}
	}

	if decreasing && increasing {
		return false
	}

	return true
}

// absDif calculates the absDif between to integers
func absDif(a, b int) int {
	if a-b > 0 {
		return a - b
	}
	return b - a
}
