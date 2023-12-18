package main

import (
	"fmt"
)

type Node struct {
	data int
	next *Node
}

type LinkedList struct {
	head *Node
}

// Insert in the beginning of the linked list
func (l *LinkedList) Insert(data int) {
	n := Node{
		data: data,
		next: nil,
	}

	if l.head == nil {
		l.head = &n
	} else {
		n.next = l.head
		l.head = &n
	}
}

func (l *LinkedList) Display() {
	if l.head == nil {
		fmt.Println("linked list is empty")
		return
	}

	cur := l.head
	for cur != nil {
		fmt.Printf("%d-> ", cur.data)
		cur = cur.next
	}
}

// Remove all ocurrencies of data
func (l *LinkedList) Remove(data int) {
	if l.head == nil {
		fmt.Println("linked list is empty")
		return
	}

	prev := l.head
	cur := l.head.next

	for cur != nil {
		if cur.data == data {
			prev.next = cur.next
		}
		prev = cur
		cur = cur.next
	}
}

func main() {
	l := LinkedList{
		head: nil,
	}

	l.Insert(1)
	l.Insert(2)
	l.Insert(4)
	l.Insert(2)
	l.Insert(3)
	l.Display()
	l.Remove(2)
	l.Display()
}
