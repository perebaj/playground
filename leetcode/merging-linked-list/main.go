// https://www.algoexpert.io/questions/merging-linked-lists

package main

// This is an input struct. Do not edit.
type LinkedList struct {
	Value int
	Next  *LinkedList
}

// o(n²)
func MergingLinkedLists(linkedListOne *LinkedList, linkedListTwo *LinkedList) *LinkedList {
	dummyListOne := linkedListOne
	dummyListTwo := linkedListTwo
	// n^2
	for dummyListOne != nil {
		dummyListTwo = linkedListTwo
		for dummyListTwo != nil {
			if dummyListOne.Value == dummyListTwo.Value {
				return dummyListTwo
			}
			dummyListTwo = dummyListTwo.Next
		}
		dummyListOne = dummyListOne.Next
	}

	return nil
}

// it's possible to use a HM and event the magic solution
// o(n+m) space: o(1)
func MergingLinkedLists2(linkedListOne *LinkedList, linkedListTwo *LinkedList) *LinkedList {
	dummyListOne := linkedListOne
	dummyListTwo := linkedListTwo
	// n^2
	for dummyListOne != nil {
		dummyListTwo = linkedListTwo
		for dummyListTwo != nil {
			if dummyListOne.Value == dummyListTwo.Value {
				return dummyListTwo
			}
			dummyListTwo = dummyListTwo.Next
		}
		dummyListOne = dummyListOne.Next
	}

	return nil
}
