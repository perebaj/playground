package adventofcode24

import "fmt"

func Day2() {
	leftList, rightList := ReadList("./day2.txt")

	rightMap := make(map[int]int)

	for _, v := range rightList {
		rightMap[v]++
	}

	var result int
	for _, v := range leftList {
		rightMapVal, ok := rightMap[v]
		if ok {
			fmt.Printf("left value %d exists in the right map x: %d times\n", v, rightMapVal)
			result = result + (v * rightMapVal)
		}
	}

	fmt.Println(result)
}
