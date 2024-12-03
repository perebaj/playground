package adventofcode24

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func Day2() {
	elements := Read("./day2tmp.txt")
	// fmt.Println(elements)

	var results int
	for _, v := range elements {
		fmt.Printf("evaluating the array %v\n", v)
		if !isMonotonic(v) {
			for i := 0; i < len(v); i++ {
				tmp := make([]string, len(v))
				copy(tmp, v)
				aux := slices.Delete(tmp, i, i+1)
				fmt.Printf("evaluating the array %v\n", aux)
				if isMonotonic(aux) {
					fmt.Println("is monotonic")
					results++
					break
				}
			}
		} else {
			fmt.Println("is monotonic")
			results++
		}
	}

	fmt.Println(results)
}

func Read(path string) [][]string {
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

func absDif(a, b int) int {
	if a-b > 0 {
		return a - b
	}
	return b - a
}
