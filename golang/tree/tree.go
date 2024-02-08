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

func (t *Tree) InsertLeft(value int) {
	if t.Root.Left == nil {
		t.Root.Left = &Node{
			Value: value,
		}
	} else {
		newLeftNode := Node{
			Value: value,
		}
		newLeftNode.Left = t.Root.Left // newLeftNode.Left -> T.Root.Left
		t.Root.Left = &newLeftNode     // t.Root.Left -> newLeftNode
	}
}

func (t *Tree) InsertRight(value int) {
	if t.Root.Right == nil {
		t.Root.Right = &Node{
			Value: value,
		}
	} else {
		newRigtNode := Node{
			Value: value,
		}
		newRigtNode.Right = t.Root.Right
		t.Root.Right = &newRigtNode
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

	fmt.Println("Inserting Left side")
	fmt.Println(t.Root.Value) // 1
	t.InsertLeft(2)
	fmt.Println(t.Root.Left.Value) // 2
	t.InsertLeft(4)
	fmt.Println(t.Root.Left.Value)      // 4
	fmt.Println(t.Root.Left.Left.Value) // 2

	fmt.Println("Inserting Right side")
	t.InsertRight(3)
	fmt.Println(t.Root.Value)       // 1
	fmt.Println(t.Root.Right.Value) // 3
	t.InsertRight(2)
	fmt.Println(t.Root.Right.Value)       // 2
	fmt.Println(t.Root.Right.Right.Value) // 3
}
