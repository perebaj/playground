package main

import "fmt"

type Node struct {
	value      string
	prev, next *Node
}

type DoublyLinkedList struct {
	head *Node
	tail *Node
}

func (b *BrowserHistory) AddNode(page string) {
	node := Node{
		value: page,
	}
	node.next = b.Chain.head.next
	node.prev = b.Chain.head
	b.Chain.head.next.prev = &node
	b.Chain.head.next = &node
}

func (b *BrowserHistory) Traverse() {
	current := b.Chain.head
	for current != nil {
		fmt.Println("traverse", current.value)
		current = current.next
	}
}

type BrowserHistory struct {
	Chain       *DoublyLinkedList
	CurPosition map[string]*Node
}

func Constructor(homepage string) BrowserHistory {
	head := &Node{}
	tail := &Node{}
	head.next = tail
	head.prev = nil
	tail.prev = head
	tail.next = nil

	return BrowserHistory{
		Chain: &DoublyLinkedList{
			head: head,
			tail: tail,
		},
		CurPosition: make(map[string]*Node),
	}
}

// func (b *BrowserHistory) Visit(url string) {

// }

// func (b *BrowserHistory) Back(steps int) string {

// }

// func (b *BrowserHistory) Forward(steps int) string {

// }

func main() {
	n := Constructor("teste")

	n.AddNode("jojo.com")
	n.AddNode("ondehj.com")

	fmt.Println(n.Chain.head.next.value)
	fmt.Println(n.Chain.head.next.next.value)
	n.Traverse()
}
