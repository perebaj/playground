package main

import (
	"errors"
	"io"
	"os"
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
