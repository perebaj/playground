package main

import "fmt"

type Node struct {
	key, value int
	prev, next *Node
}

type LRUCache struct {
	capacity int
	cache    map[int]*Node
	head     *Node
	tail     *Node
}

func Constructor(capacity int) LRUCache {
	head := &Node{}
	tail := &Node{}
	head.next = tail
	head.prev = nil
	tail.prev = head
	tail.next = nil
	return LRUCache{
		capacity: capacity,
		cache:    map[int]*Node{},
		head:     head,
		tail:     tail,
	}
}

func (l *LRUCache) DeleteNode(node *Node) {
	prev := node.prev
	next := node.next
	prev.next = next
	next.prev = prev

	//delete from hashmap
	delete(l.cache, node.key)
}

func (l *LRUCache) AddNode(node *Node) {
	node.next = l.head.next
	node.prev = l.head
	l.head.next.prev = node
	l.head.next = node
	//add in Hashmap
	l.cache[node.key] = node
}

func (l *LRUCache) Traverse() {
	current := l.head
	for current != nil {
		fmt.Println("traverse", current.key, current.value)
		current = current.next
	}
}

func (l *LRUCache) isMaxCapacity() bool {
	if len(l.cache) == l.capacity {
		return true
	} else {
		return false
	}
}

func (l *LRUCache) Get(key int) int {
	v, ok := l.cache[key]
	if ok {
		l.DeleteNode(v)
		l.AddNode(v)
		return v.value
	} else {
		return -1
	}
}

func (l *LRUCache) Put(key, value int) {
	v, ok := l.cache[key]
	if ok {
		l.DeleteNode(v)
		// newNode := Node{
		// 	key:   key,
		// 	value: value,
		// }
		// l.AddNode(&newNode)
	} else if l.isMaxCapacity() {
		//delete last node
		lastNode := l.tail.prev
		l.DeleteNode(lastNode)
	}
	newNode := Node{
		key:   key,
		value: value,
	}
	l.AddNode(&newNode)

}

func main() {
	//[[2],[2],[2,6],[1],[1,5],[1,2],[1],[2]]

	n := Constructor(2)
	fmt.Println(n.Get(2))
	n.Put(2, 6)
	fmt.Println(n.Get(1))
	n.Put(1, 5)
	n.Traverse()
	n.Put(1, 2)
	n.Traverse()
	fmt.Println(n.Get(1))
	fmt.Println(n.Get(2))
	n.Traverse()
}
