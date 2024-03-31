package main

func main() {
	nums := []int{-1, 2, -3, 4} // 2 -1 4 -3
	rearrangeArray(nums)
	nums = []int{1, -22, 3, -4, -1, 4} // 1 -22 3 -4 4 -1
	rearrangeArray(nums)

}

func rearrangeArray(nums []int) []int {
	var negative []int
	var positive []int
	for _, v := range nums {
		if v > 0 {
			positive = append(positive, v)
		} else {
			negative = append(negative, v)
		}
	}
	var resp []int
	var pointer int
	for i := 0; i < len(positive)*2; i++ {
		if i%2 == 0 { // even numbers
			resp = append(resp, positive[pointer])
		} else { //odd numbers
			resp = append(resp, negative[pointer])
			pointer++
		}
	}

	return resp
}
