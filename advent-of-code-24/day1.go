package adventofcode24

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
)

func Day1() {
	leftList, rightList := ReadList("./day1.txt")

	sort.Ints(leftList)
	sort.Ints(rightList)

	var result int
	for i := 0; i < len(rightList); i++ {
		abs := absDistance(rightList[i], leftList[i])
		result = result + abs
	}

	fmt.Println(result)
}

func ReadList(pathFile string) ([]int, []int) {
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
		fmt.Errorf("%v", err)
	}

	return leftList, rightList
}

func absDistance(a, b int) int {
	if a-b < 0 {
		return b - a
	}
	return a - b
}
