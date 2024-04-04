package main

type ListNode struct {
	Val  int
	Next *ListNode
}

func AddNode(head *ListNode, val int) {
	for head.Next != nil {
		head = head.Next
	}
	head.Next = &ListNode{Val: val}
}

func traverseLL(head *ListNode) {
	for head != nil {
		println(head.Val)
		head = head.Next
	}
}

func main() {
	head := &ListNode{Val: 1}
	AddNode(head, 2)
	AddNode(head, 3)
	AddNode(head, 4)
	AddNode(head, 5)
	reorderList(head) // 1 -> 5 -> 2 -> 4 -> 3
}

func reorderList(head *ListNode) {
	var nodeList []*ListNode
	for head != nil {
		nodeList = append(nodeList, head)
		head = head.Next
	}

	newTailIndex := (len(nodeList) / 2)

	start, end := 0, len(nodeList)-1
	for start < end {
		nodeList[start].Next = nodeList[end]
		// we need to made the last pointer point to nil, otherwise, we will have a cycle in the linked list.
		start++
		nodeList[end].Next = nodeList[start]
		end--
	}
	nodeList[newTailIndex].Next = nil
	head = nodeList[0]
	// traverseLL(head)
}
