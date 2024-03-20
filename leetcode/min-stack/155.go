package main

import (
	"fmt"
)

type MinStack struct {
	Stack    []int
	MinStack []int
}

func Constructor() MinStack {
	return MinStack{
		Stack:    []int{},
		MinStack: []int{},
	}
}

func (m *MinStack) Push(val int) {
	m.Stack = append(m.Stack, val)
	if len(m.MinStack) > 0 {
		minimum := min(val, m.MinStack[len(m.MinStack)-1])
		m.MinStack = append(m.MinStack, minimum)
	} else {
		m.MinStack = append(m.MinStack, val)
	}
}

func (m *MinStack) Pop() {
	m.Stack = m.Stack[:len(m.Stack)-1]
	m.MinStack = m.MinStack[:len(m.MinStack)-1]
}

func (m *MinStack) Top() int {
	return m.Stack[len(m.Stack)-1]
}

func (m *MinStack) GetMin() int {
	min := m.MinStack[len(m.MinStack)-1]
	return min
}

/**
The tricky to solve problem is create a smart way to store the minimum values

Basiclly we can create a

to normal stack will follow the push and pop operations normally, but the minimumStack should follow up a different approach

each time that we are inserting a new element we need to aswer if this new element is smaller then the current one, at this way


1) 1   1
2) 0   0
3) 1   0

obj.pop() // removing the one and also the minimum element related to one(that is zero)


*/

func main() {
	obj := Constructor()
	obj.Push(-2)
	obj.Push(0)
	obj.Push(-3)
	fmt.Println(obj.GetMin())
	obj.Pop()
	fmt.Println(obj.Top())
	// fmt.Println(obj.GetMin())

}
