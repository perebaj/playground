package main

import "fmt"

// ll1: 2 -> 4 -> 7 -> 1
// ll2: 9 -> 4 -> 5

// This is an input struct. Do not edit.
type LinkedList struct {
	Value int
	Next  *LinkedList
}

func SumOfLinkedLists(linkedListOne *LinkedList, linkedListTwo *LinkedList) *LinkedList {
	var resp LinkedList
	dummy := &resp
	var l1Val int
	var l2Val int
	var carry int
	for linkedListOne != nil || linkedListTwo != nil || carry != 0 {
		if linkedListOne == nil {
			l1Val = 0
		} else {
			l1Val = linkedListOne.Value
		}

		if linkedListTwo == nil {
			l2Val = 0
		} else {
			l2Val = linkedListTwo.Value
		}
		fmt.Println(l1Val, l2Val, carry)
		auxSum := (l1Val + l2Val + carry) % 10
		carry = (l1Val + l2Val + carry) / 10
		fmt.Println(auxSum, carry)
		dummy.Next = &LinkedList{
			Value: auxSum,
		}

		dummy = dummy.Next
		if linkedListOne != nil {
			linkedListOne = linkedListOne.Next
		}

		if linkedListTwo != nil {
			linkedListTwo = linkedListTwo.Next
		}
	}

	return resp.Next
}
