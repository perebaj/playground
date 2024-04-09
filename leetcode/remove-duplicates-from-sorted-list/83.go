package main

type ListNode struct {
	Val  int
	Next *ListNode
}

func deleteDuplicates(head *ListNode) *ListNode {
	if head == nil {
		return head
	}
	prev := head
	dummy := prev
	cur := head.Next
	for cur != nil {
		if prev.Val == cur.Val {
			prev.Next = cur.Next
			cur = prev.Next
		} else {
			prev = cur
			cur = cur.Next
		}
	}

	return dummy
}
