package main

import "fmt"

type Cashier struct {
	NCustumer     int
	Discount      int
	ProduPrice    map[int]int
	CustomerCount int
}

func Constructor(n int, discount int, products []int, prices []int) Cashier {
	m := make(map[int]int)

	for i := 0; i < len(products); i++ {
		m[products[i]] = prices[i]
	}

	return Cashier{
		ProduPrice: m,
		NCustumer:  n,
		Discount:   discount,
	}
}

func (c *Cashier) GetBill(product []int, amount []int) float64 {
	c.CustomerCount++
	var bill int
	for i := 0; i < len(product); i++ {
		price := c.ProduPrice[product[i]]
		bill += price * amount[i]
	}
	var newBill float64
	if c.CustomerCount%c.NCustumer == 0 {
		//apply discoutn
		// fmt.Println("asdlkj")
		newBill = float64(bill) * ((100 - float64(c.Discount)) / 100)
	}
	if newBill != 0 {
		return newBill
	}
	return float64(bill)
}

func main() {
	n := Constructor(3, 50, []int{1, 2, 3, 4, 5, 6, 7}, []int{100, 200, 300, 400, 300, 200, 100})
	fmt.Println(n.GetBill([]int{1, 2}, []int{1, 2}))                               // return 500.0
	fmt.Println(n.GetBill([]int{3, 7}, []int{10, 10}))                             // return 4000.0
	fmt.Println(n.GetBill([]int{1, 2, 3, 4, 5, 6, 7}, []int{1, 1, 1, 1, 1, 1, 1})) // return 800.0
}
