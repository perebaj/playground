package main

import (
	"reflect"
	"testing"
)

func TestMain(t *testing.T) {
	input1 := []int{2, 7, 11, 15}
	target1 := 9
	got := resolution(input1, target1)
	want := []int{0, 1}
	if !reflect.DeepEqual(want, got) {
		t.Error("expeting a diffent value")
	}
}
