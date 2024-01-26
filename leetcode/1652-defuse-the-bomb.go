package main

import (
	"fmt"
)

func main() {
	// code := []int{5, 7, 1, 4}
	// k := 3
	// fmt.Println(decrypt(code, k)) // [12,10,16,13]

	// code = []int{1, 2, 3, 4}
	// k = 0
	// fmt.Println(decrypt(code, k)) // [0,0,0,0]

	// code := []int{2, 4, 9, 3} //12,5,6,13
	// k := -2
	// fmt.Println(decrypt(code, k))

	// code := []int{10, 5, 7, 7, 3, 2, 10, 3, 6, 9, 1, 6}
	// k := -4
	// fmt.Println(decrypt(code, k))

	code := []int{2, 4, 9, 3, 2, 4, 9, 3} //size8
	n := len(code)
	k := 3
	for i := 0; i < len(code); i++ {
		for j := 0; j < k; j++ {
			fmt.Println(n - j - i - 1)
			// fmt.Println((-1 - i) - n)

		}
		fmt.Println("dif loop")
	}

}

func decrypt(code []int, k int) []int {
	n := len(code)
	var sum int
	res := make([]int, n)

	for i := 0; i < n; i++ {
		if k > 0 {
			for j := 0; j < k; j++ {
				sum = sum + code[(i+j+1)%n]
			}
			res[i] = sum
			sum = 0
		}
		// THe problem is here
		if k < 0 {
			fmt.Printf("index %d\n", i)
			for j := 0; j < (-1)*k; j++ {
				if i >= (-1)*k {
					index := (i - 1 - j)
					fmt.Printf("case1 | index %d\n", index)
					sum = sum + code[index]
					// fmt.Printf("sum:%d\n", sum)
				} else {
					index := (n - j - 1 - i) % n
					fmt.Printf("case2 | index %d | valores i=%d, j=%d, n=%d \n", index, i, j, n)
					sum = sum + code[index]
				}
				// fmt.Printf("sum:%d\n", sum)
			}
			res[i] = sum
			sum = 0
		}
	}

	return res
}

/*   0,1,2,3,4,5,6,7
[2,4,9,3,2,4,9,3] k=-3

index0 = 7,6,5
index1=0,7,6
index2=1,0,7
index3=2,1,0
index4=3,2,1
index5=4,3,2
index6=5,4,3
index7=6,5,4


for i->len(code) {
    for j 0->k{
        n-1-j
    }
}
*/
