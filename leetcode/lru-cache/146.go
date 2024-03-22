package main

import "fmt"

type DoublyLinkedList struct {
	Head *Node
}

type Node struct {
	Data int
	Key  int
	Next *Node
	Prev *Node
}

type LRUCache struct {
	DoublyLL DoublyLinkedList
	HM       map[int]int
	Capacity int
}

func Constructor(capacity int) LRUCache {
	return LRUCache{
		DoublyLL: DoublyLinkedList{},
		HM:       make(map[int]int),
		Capacity: capacity,
	}
}

func (l *LRUCache) isMaxCapacity(key int) bool {
	_, ok := l.HM[key]
	if ok {
		return false
	} else {
		return len(l.HM) == l.Capacity
	}
}

func (n *DoublyLinkedList) DeleteEnd() int {
	if n.Head == nil {
		return -1
	} else {
		aux := n.Head
		for {
			if aux.Next == nil && aux.Prev != nil { // Deleting the last node
				prevNode := aux.Prev
				prevNode.Next = nil
				return aux.Key
			} else if aux.Next == nil && aux.Prev == nil {
				n.Head = nil
				return aux.Key
			} else {
				aux = aux.Next
			}
		}
	}

}

// WasAccessed change
func (n *DoublyLinkedList) WasAccessed(key int) {
	if n.Head == nil {
		return
	} else {
		aux := n.Head
		for aux != nil {
			if aux.Key == key {
				prevNode := aux.Prev
				nextNode := aux.Next
				if prevNode == nil && nextNode == nil || prevNode == nil {
					break
				} else {
					if nextNode != nil {
						nextNode.Prev = prevNode
					}
					if prevNode != nil {
						prevNode.Next = nextNode
					}
					aux.Prev = nil
					aux.Next = n.Head
					n.Head.Prev = aux
					n.Head = aux
					return
				}
			}

			aux = aux.Next
		}
	}
}

func (n *DoublyLinkedList) InsertBegin(data int, key int) {
	newNode := Node{
		Data: data,
		Prev: nil,
		Key:  key,
	}
	// Verify if the key already exists, if so, update the value
	aux := n.Head
	for {
		if aux == nil {
			break
		} else {
			if aux.Key == key {
				aux.Data = data
				n.WasAccessed(key)
				return
			}
			aux = aux.Next
		}
	}

	if n.Head == nil {
		newNode.Next = nil
		newNode.Prev = nil
	} else {
		n.Head.Prev = &newNode
		newNode.Next = n.Head
	}
	n.Head = &newNode
}

func (n *DoublyLinkedList) traverse() {
	aux := n.Head
	for {
		if aux == nil {
			break
		} else {
			aux = aux.Next
		}
	}

}

func (l *LRUCache) Put(key int, value int) {
	if l.isMaxCapacity(key) {
		deletedKey := l.DoublyLL.DeleteEnd()
		if deletedKey == -1 {
			return
		} else {
			delete(l.HM, deletedKey)
		}
		// l.DoublyLL.InsertBegin(value, key)
		// l.HM[key] = value
	}
	l.DoublyLL.InsertBegin(value, key)
	l.HM[key] = value
}

func (l *LRUCache) Get(key int) int {
	v, ok := l.HM[key]
	if !ok {
		return -1
	} else {
		l.DoublyLL.WasAccessed(key)
		return v
	}
}

func main() {
	//[[3],[1,1],[2,2],[3,3],[4,4],[4],[3],[2],[1],[5,5],[1],[2],[3],[4],[5]]

	LRUCache := Constructor(3)
	LRUCache.Put(1, 1)
	LRUCache.Put(2, 2)
	LRUCache.Put(3, 3)
	LRUCache.Put(4, 4)
	fmt.Println(LRUCache.Get(4))
	fmt.Println(LRUCache.Get(3))
	fmt.Println(LRUCache.Get(2))
	LRUCache.DoublyLL.traverse()
	fmt.Println(LRUCache.Get(1))
	LRUCache.Put(5, 5)
	fmt.Println(LRUCache.HM)
	fmt.Println(LRUCache.Get(1))
	fmt.Println(LRUCache.Get(2))
	fmt.Println(LRUCache.Get(3))
	fmt.Println(LRUCache.Get(4))
	fmt.Println(LRUCache.Get(5))
}

/**
 * Your LRUCache object will be instantiated and called as such:
 * obj := Constructor(capacity);
 * param_1 := obj.Get(key);
 * obj.Put(key,value);
 */
