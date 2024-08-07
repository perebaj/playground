package main

import (
	"errors"
	"io"
	"os"
	"slices"
	"testing"
)

var ErrInvalidNode = errors.New("invalid node")

func TestInsertLeft(t *testing.T) {
	tree := &Tree{
		Root: &Node{
			Value: 1,
		},
	}

	if tree.Root.Left != nil {
		t.Error(ErrInvalidNode)
	}

	if tree.Root.Right != nil {
		t.Error(ErrInvalidNode)
	}

	if tree.Root.Value != 1 {
		t.Error(ErrInvalidNode)
	}

	tree.Root.InsertLeft(3)

	if tree.Root.Left.Value != 3 {
		t.Error(ErrInvalidNode)
	}

	tree.Root.InsertLeft(2)

	if tree.Root.Left.Value != 2 {
		t.Error(ErrInvalidNode)
	}

	if tree.Root.Left.Left.Value != 3 {
		t.Error(ErrInvalidNode)
	}
}

func TestInsertRight(t *testing.T) {
	tree := &Tree{
		Root: &Node{
			Value: 1,
		},
	}

	if tree.Root.Left != nil {
		t.Error(ErrInvalidNode)
	}

	if tree.Root.Right != nil {
		t.Error(ErrInvalidNode)
	}

	if tree.Root.Value != 1 {
		t.Error(ErrInvalidNode)
	}

	tree.Root.InsertRight(3)

	if tree.Root.Right.Value != 3 {
		t.Error(ErrInvalidNode)
	}

	tree.Root.InsertRight(2)

	if tree.Root.Right.Value != 2 {
		t.Error(ErrInvalidNode)
	}

	if tree.Root.Right.Right.Value != 3 {
		t.Error(ErrInvalidNode)
	}
}

func TestInorder(t *testing.T) {
	tree := &Tree{
		Root: &Node{
			Value: 1,
		},
	}

	tree.Root.InsertLeft(2)
	tree.Root.InsertRight(5)

	node2 := tree.Root.Left
	node5 := tree.Root.Right

	node2.InsertLeft(3)
	node2.InsertRight(4)

	node5.InsertLeft(6)
	node5.InsertRight(7)

	expected := "3-2-4-1-6-5-7-"
	storeStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	inorder(tree.Root)

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = storeStdout

	if string(out) != expected {
		t.Errorf("Expected %s, got %s", expected, out)
	}
}

func TestPostorder(t *testing.T) {
	tree := &Tree{
		Root: &Node{
			Value: 1,
		},
	}

	tree.Root.InsertLeft(2)
	tree.Root.InsertRight(5)

	node2 := tree.Root.Left
	node5 := tree.Root.Right

	node2.InsertLeft(3)
	node2.InsertRight(4)

	node5.InsertLeft(6)
	node5.InsertRight(7)

	expected := "3-4-2-6-7-5-1-"
	storeStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	postoder(tree.Root)

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = storeStdout

	if string(out) != expected {
		t.Errorf("Expected %s, got %s", expected, out)
	}
}

func TestPreorder(t *testing.T) {
	tree := &Tree{
		Root: &Node{
			Value: 1,
		},
	}

	tree.Root.InsertLeft(2)
	tree.Root.InsertRight(5)

	node2 := tree.Root.Left
	node5 := tree.Root.Right

	node2.InsertLeft(3)
	node2.InsertRight(4)

	node5.InsertLeft(6)
	node5.InsertRight(7)

	expected := "1-2-3-4-5-6-7-"
	storeStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	preoder(tree.Root)

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = storeStdout

	if string(out) != expected {
		t.Errorf("Expected %s, got %s", expected, out)
	}
}

func TestBinarySearchTreeInsert(t *testing.T) {
	tree := Tree{
		Root: &Node{
			Value: 50,
		},
	}

	want := 52
	tree.Root.BinarySeachTreeInsert(want)

	if tree.Root.Value != 50 {
		t.Error(ErrInvalidNode)
	}

	got := tree.Root.Right.Value
	if got != want {
		t.Errorf("Unexpeted value, want %d | got %d", want, got)
	}

	want = 51
	tree.Root.BinarySeachTreeInsert(51)

	got = tree.Root.Right.Left.Value

	if got != want {
		t.Errorf("Unexpeted value, want %d | got %d", want, got)
	}

	want = 50
	tree.Root.BinarySeachTreeInsert(want)

	got = tree.Root.Left.Value

	if got != want {
		t.Errorf("Unexpeted value, want %d | got %d", want, got)
	}
}

func TestValueExists(t *testing.T) {
	tree := Tree{
		Root: &Node{
			Value: 52,
		},
	}

	tree.Root.BinarySeachTreeInsert(53)
	tree.Root.BinarySeachTreeInsert(40)
	tree.Root.BinarySeachTreeInsert(35)
	tree.Root.BinarySeachTreeInsert(43)
	tree.Root.BinarySeachTreeInsert(60)

	ok := tree.Root.valueExists(40)

	if !ok {
		t.Error("Expecting that the value exists")
	}

	ok = tree.Root.valueExists(100)
	if ok {
		t.Error("Expecting that the value doens't exists")
	}

	//nil tree

	tree2 := Tree{
		Root: &Node{},
	}

	ok = tree2.Root.valueExists(10)
	if ok {
		t.Error("Passing nil Tree, expecting false")
	}
}

func TestInorderSlice(t *testing.T) {
	// temp := make([]int, 10)

	tree := &Tree{
		Root: &Node{
			Value: 1,
		},
	}

	tree.Root.InsertLeft(2)
	tree.Root.InsertRight(5)

	node2 := tree.Root.Left
	node5 := tree.Root.Right

	node2.InsertLeft(3)
	node2.InsertRight(4)

	node5.InsertLeft(6)
	node5.InsertRight(7)

	want := []int{3, 2, 4, 1, 6, 5, 7}

	got := inorderSlice(tree.Root)

	if !slices.Equal(want, got) {
		t.Errorf("want %v | got %v", want, got)
	}
}
