package main

import "fmt"

func main() {
	r := maxNumberOfBalloons("loonbalxballpoon")
	fmt.Println("resp:", r)
}

func maxNumberOfBalloons(text string) int {
	b, a, l, o, n := 0, 0, 0, 0, 0
	// b=1,a=1,l=2, o=2, n=1
	for _, r := range text {
		switch r {
		case 'b':
			b++
		case 'a':
			a++
		case 'l':
			l++
		case 'o':
			o++
		case 'n':
			n++
		}
	}
	r := min(a, b)
	r2 := min(r, n)
	r3 := min(r2, l/2)
	r4 := min(r3, o/2)

	return r4

}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
