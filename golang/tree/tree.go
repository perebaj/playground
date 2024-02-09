package main

import "fmt"

// Binary tree implementation
// Reference: https://www.freecodecamp.org/news/all-you-need-to-know-about-tree-data-structures-bceacb85490c/

type Tree struct {
	Root *Node
}

type Node struct {
	Value int
	Left  *Node
	Right *Node
}

func (n *Node) InsertLeft(value int) {
	if n.Left == nil {
		n.Left = &Node{
			Value: value,
		}
	} else {
		newLeftNode := Node{
			Value: value,
		}
		newLeftNode.Left = n.Left // newLeftNode.Left -> T.Left
		n.Left = &newLeftNode     // t.Root.Left -> newLeftNode
	}
}

func (n *Node) InsertRight(value int) {
	if n.Right == nil {
		n.Right = &Node{
			Value: value,
		}
	} else {
		newRigtNode := Node{
			Value: value,
		}
		newRigtNode.Right = n.Right
		n.Right = &newRigtNode
	}
}

/*
Depth-First traverses
It's easy to memorize the order of the traverses, using the place where you need to print the value of the node:
- In-Order: Between the left and right nodes
- Post-Order: After the left and right nodes
- Pre-Order: Before the left and right nodes
*/
func inorder(n *Node) {
	if n != nil {
		inorder(n.Left)
		fmt.Print(n.Value)
		fmt.Print("-")
		inorder(n.Right)
	}
}

func postoder(n *Node) {
	if n != nil {
		postoder(n.Left)
		postoder(n.Right)
		fmt.Print(n.Value)
		fmt.Print("-")
	}
}

func preoder(n *Node) {
	if n != nil {
		fmt.Print(n.Value)
		fmt.Print("-")
		preoder(n.Left)
		preoder(n.Right)
	}
}

func main() {
	t := &Tree{
		Root: &Node{
			Value: 1,
		},
	}

	fmt.Println(t.Root.Value) // 1
	fmt.Println(t.Root.Left)  // Nil
	fmt.Println(t.Root.Right) // Nil

	t.Root.InsertLeft(2)
	t.Root.InsertRight(5)

	node2 := t.Root.Left
	node5 := t.Root.Right

	node2.InsertLeft(3)
	node2.InsertRight(4)

	node5.InsertLeft(6)
	node5.InsertRight(7)

	// All of the following traverses are Depth-First traverses, that
	fmt.Println("In-Order traverses")
	inorder(t.Root) // Expected result: 3-2-4-1-6-5-7
	fmt.Println("\nPost-Order traverses")
	postoder(t.Root) // Expected result: 3-4-2-6-7-5-1
	fmt.Println("\nPre-Order traverses")
	preoder(t.Root) // Expected result: 1-2-3-4-5-6-7

}
