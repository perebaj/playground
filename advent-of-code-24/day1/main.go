package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
)

func main() {
	Day1Part1()
	Day1Part2()
}

func Day1Part1() {
	leftList, rightList := readInput("./day1.txt")

	sort.Ints(leftList)
	sort.Ints(rightList)

	var result int
	for i := 0; i < len(rightList); i++ {
		abs := absDistance(rightList[i], leftList[i])
		result = result + abs
	}

	fmt.Println("Result Day 1 Part 2: ", result)
}

func Day1Part2() {
	leftList, rightList := readInput("./day1.txt")

	rightMap := make(map[int]int)

	for _, v := range rightList {
		rightMap[v]++
	}

	var result int
	for _, v := range leftList {
		rightMapVal, ok := rightMap[v]
		if ok {
			result = result + (v * rightMapVal)
		}
	}

	fmt.Println("Result Day 1 Part 1: ", result)
}

func readInput(pathFile string) ([]int, []int) {
	file, err := os.Open(pathFile)
	if err != nil {
		panic("failed to read")
	}

	defer file.Close()

	var leftList, rightList []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		r, _ := regexp.Compile(`(?m)(\d.*)   (\d.*)`)
		all := r.FindStringSubmatch(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		leftAux, _ := strconv.Atoi(all[1])
		rightAux, _ := strconv.Atoi(all[2])
		leftList = append(leftList, leftAux)
		rightList = append(rightList, rightAux)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return leftList, rightList
}

// Calculate the absolute difference between 2 integers
func absDistance(a, b int) int {
	if a-b < 0 {
		return b - a
	}
	return a - b
}
