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
	day1()

	// var re = regexp.MustCompile(`(?m)(\d.*)   (\d.*)`)
	// var str = `40758   64262`

	// match := re.FindStringSubmatch(str)
	// fmt.Println(match, len(match))
	// fmt.Println(match[1], match[2])

}

func day1() {
	file, err := os.Open("./day1.txt")
	if err != nil {
		panic("failed to read")
	}

	defer file.Close()

	var leftList, rightList []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// fmt.Println(scanner.Text())
		r, _ := regexp.Compile(`(?m)(\d.*)   (\d.*)`)
		all := r.FindStringSubmatch(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		//convert str to int32
		leftAux, _ := strconv.Atoi(all[1])
		rightAux, _ := strconv.Atoi(all[2])
		leftList = append(leftList, leftAux)
		rightList = append(rightList, rightAux)

	}

	if err := scanner.Err(); err != nil {
		fmt.Errorf("%v", err)
	}

	// fmt.Println(leftList[:2], rightList[:2])
	sort.Ints(leftList)
	sort.Ints(rightList)
	fmt.Println(len(leftList), leftList[:10])
	fmt.Println(len(rightList), rightList[:10])

	var result int
	for i := 0; i < len(rightList); i++ {
		abs := absDistance(rightList[i], leftList[i])
		fmt.Println(abs)
		result = result + abs
	}

	fmt.Println(result)
}

func absDistance(a, b int) int {
	if a-b < 0 {
		return b - a
	}
	return a - b
}

// func main() {
//     file, err := os.Open("/path/to/file.txt")
//     if err != nil {
//         log.Fatal(err)
//     }
//     defer file.Close()

//     scanner := bufio.NewScanner(file)
//     // optionally, resize scanner's capacity for lines over 64K, see next example
//     for scanner.Scan() {
//         fmt.Println(scanner.Text())
//     }

//     if err := scanner.Err(); err != nil {
//         log.Fatal(err)
//     }
// }
