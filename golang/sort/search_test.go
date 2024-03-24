package main

import "testing"

func Test(t *testing.T) {
	sortSearch()
}

func TestIndex(t *testing.T) {
	index := Index()
	if index != 0 {
		t.Errorf("Index() = %d; want 0", index)
	}
}

func TestContructor(t *testing.T) {
	Constructor([]int{1, 2, 3, 2, 5, 2, 7, 8, 9})
}
