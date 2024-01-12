package main

import (
	"fmt"
	"reflect"
)

func main() {
	path := "SN"

	fmt.Println(isPathCrossing(path))
}

func isPathCrossing(path string) bool {

	cordi := []int{0, 0}
	var auxSlice [][]int
	auxSlice = append(auxSlice, []int{0, 0})
	for i := 0; i < len(path); i++ {
		auxCord := make([]int, 2)
		switch string(path[i]) {
		case "N":
			cordi[1] = cordi[1] + 1
		case "W":
			cordi[0] = cordi[0] - 1
		case "S":
			cordi[1] = cordi[1] - 1
		case "E":
			cordi[0] = cordi[0] + 1
		}
		copy(auxCord, cordi)
		auxSlice = append(auxSlice, auxCord)
	}

	for i := 0; i < len(auxSlice); i++ {
		for j := i + 1; j < len(auxSlice); j++ {
			if reflect.DeepEqual(auxSlice[i], auxSlice[j]) {
				return true
			}
		}
	}

	return false
}

/* simplified
func isPathCrossing(path string) bool {
    visited := make(map[string] bool)

    visited["00"] = true

    x := 0
    y := 0

    for _, val := range path {
        switch val {
        case 'N':
            y++
        case 'S':
            y--
        case 'E':
            x++
        default:
            x--
        }


        key := fmt.Sprintf("%d%d", x, y)

        _, ok := visited[key]

        if ok {
            return true
        }

        visited[key] = true
    }

    return false
}
*/
