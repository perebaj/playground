package main

import (
	"errors"
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

	tree.InsertLeft(3)

	if tree.Root.Left.Value != 3 {
		t.Error(ErrInvalidNode)
	}

	tree.InsertLeft(2)

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

	tree.InsertRight(3)

	if tree.Root.Right.Value != 3 {
		t.Error(ErrInvalidNode)
	}

	tree.InsertRight(2)

	if tree.Root.Right.Value != 2 {
		t.Error(ErrInvalidNode)
	}

	if tree.Root.Right.Right.Value != 3 {
		t.Error(ErrInvalidNode)
	}
}
