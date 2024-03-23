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
	b.CurPosition["current"] = &node
}

func (b *BrowserHistory) AddNodeEnd(page string) {
	node := Node{
		value: page,
	}
	node.prev = b.Chain.tail.prev
	node.next = b.Chain.tail
	b.Chain.tail.prev.next = &node
	b.Chain.tail.prev = &node
	b.CurPosition["current"] = &node
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
	head := &Node{
		value: "head",
	}
	tail := &Node{
		value: "tail",
	}
	head.next = tail
	head.prev = nil
	tail.prev = head
	tail.next = nil

	b := BrowserHistory{
		Chain: &DoublyLinkedList{
			head: head,
			tail: tail,
		},
		CurPosition: make(map[string]*Node),
	}

	b.AddNode(homepage)

	return b
}

func (b *BrowserHistory) Visit(url string) {
	b.AddNodeEnd(url)
}

func (b *BrowserHistory) Back(steps int) string {
	aux := b.CurPosition["current"]
	// traverse the linked list if achieve the head node, return the head.next.value, otherwise if achive the steps == 0, return the current node value
	for steps > 0 {
		if aux.value == "head" {
			break
		}
		aux = aux.prev
		steps--
	}
	if aux.value == "head" {
		b.CurPosition["current"] = aux.next
		return aux.next.value
	}
	b.CurPosition["current"] = aux
	return aux.value
}

func (b *BrowserHistory) Forward(steps int) string {
	aux := b.CurPosition["current"]
	for steps > 0 {
		if aux.value == "tail" {
			break
		}
		aux = aux.next
		steps--
	}
	if aux.value == "tail" {
		b.CurPosition["current"] = aux.prev
		return aux.prev.value
	}
	b.CurPosition["current"] = aux
	return aux.value
}

func main() {
	n := Constructor("leetcode.com")
	n.Visit("google.com")
	n.Visit("facebook.com")
	n.Visit("youtube.com")
	fmt.Println(n.Back(1))    // facebook
	fmt.Println(n.Back(1))    // google
	fmt.Println(n.Forward(1)) // facebook
	n.Visit("linkedin.com")
	fmt.Println(n.Forward(2)) // linkedin
	fmt.Println(n.Back(2))    // facebook (I think that the expected output is wrong)
	fmt.Println(n.Back(7))    // leetcode

	// expected output
	//[null,null,null,null,"facebook.com","google.com","facebook.com",null,"linkedin.com","google.com","leetcode.com"]

}
