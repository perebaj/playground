package main

func main() {
	fib(2) // 1
	fib(3) // 2
	fib(4) // 3

}

func fib(n int) int {
	var fibArray []int
	for i := 0; i <= n; i++ {
		if i == 0 {
			fibArray = append(fibArray, 0)
		} else if i == 1 {
			fibArray = append(fibArray, 1)
		} else {
			res := fibArray[i-1] + fibArray[i-2]
			fibArray = append(fibArray, res)
		}
	}
	return fibArray[n]
}
