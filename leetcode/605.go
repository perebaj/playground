package main

import "fmt"

func main() {
	flowerbed := []int{1}
	n := 1
	fmt.Println(canPlaceFlowers(flowerbed, n))
}

func canPlaceFlowers(flowerbed []int, n int) bool {

	if n == 0 {
		return true
	}

	if len(flowerbed) == 1 {
		if flowerbed[0] == 0 && n == 1 {
			return true
		} else if flowerbed[0] == 1 && n == 1 {
			return false
		}
	}

	for index := 0; index < len(flowerbed); index++ {
		if index == 0 || index == len(flowerbed)-1 {
			if index == 0 && flowerbed[index] == 0 && flowerbed[index+1] != 1 {
				flowerbed[index] = 1
				n--
			}
			if index == len(flowerbed)-1 && flowerbed[index-1] != 1 && flowerbed[index] == 0 {
				flowerbed[index] = 1
				n--
			}
		} else {
			if flowerbed[index] == 0 && flowerbed[index+1] != 1 && flowerbed[index-1] != 1 {
				flowerbed[index] = 1
				n--
			}
		}
	}
	if n > 0 {
		return false
	}
	return true
}

/*
flowerbed = [1,0,0,0,1], n = 1

- put the plant in the right place = 0
- decrease 1 in the n variable

corner case len(flowerbed) <= 2

[0,1]
when index was 1, I will have an error!

thinking when len > 2

corner case 1:
	when index == 0 AND flowerbed[index] == 0 AND flowebed[index+1] != 1 {
		flowerbed[index] = 1
		n--
	}
corner case 2:
	when index == (len(flowerbed) - 1) AND flowerbed[index] == 0 AND flowerbed[index-1] != 1 {
		flowerbed[index] = 1
		n--
	}
normal case 3:
	flowerbed[index] == 0 AND flowerbed[index + 1] != 1 AND flowerbed[index -1] != 1



*/
