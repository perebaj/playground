package main

type ListNode struct {
	Val  int
	Next *ListNode
}

// The tricky of this exercise is undestand that we can sum the number in the reserved order, what is not to explicit in the exercise description
// after that we just need to know how to calculate the carry and value in a sum operation
// and define the right delimiters to stop the infinit loop
// I consider this one as a medium plus
func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	var dummy ListNode
	aux := &dummy

	var carry int
	l1Val, l2Val := 0, 0
	for l1 != nil || l2 != nil || carry > 0 {
		if l2 == nil {
			l2Val = 0
		} else {
			l2Val = l2.Val
		}

		if l1 == nil {
			l1Val = 0
		} else {
			l1Val = l1.Val
		}

		auxVal := l2Val + l1Val + carry
		carry = auxVal / 10
		auxVal = auxVal % 10
		aux.Next = &ListNode{
			Val: auxVal,
		}
		aux = aux.Next
		if l1 != nil {
			l1 = l1.Next
		}

		if l2 != nil {
			l2 = l2.Next
		}
	}
	return dummy.Next
}
