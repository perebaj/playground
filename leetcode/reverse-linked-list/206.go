package main

type ListNode struct {
	Val  int
	Next *ListNode
}

func reverseList(head *ListNode) *ListNode {
	if head == nil {
		return head
	}
	var prev *ListNode
	cur := head
	for cur != nil {
		aux := cur.Next
		cur.Next = prev
		prev = cur
		cur = aux
	}

	return prev
}
