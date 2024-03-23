package integercontainer

import "sort"

type IntegerContainerImpl struct {
	AbstractIntegerContainer
	
	Slice []int
}

func NewIntegerContainerImpl() *IntegerContainerImpl {
	IntegerContainer := IntegerContainerImpl{AbstractIntegerContainer{}, []int{}}
	return &IntegerContainer
}

func (i *IntegerContainerImpl) Add(num int) int {
	i.Slice = append(i.Slice, num)
	return len(i.Slice)
}

func (i *IntegerContainerImpl) Delete(num int) bool {
	index := -1
	for k, v := range i.Slice {
		if num == v {
			index = k
			break
		}
	}
	if index == -1 {
		return false
	}

	i.Slice = append(i.Slice[:index], i.Slice[index+1:]...)
	return true
}

func (i *IntegerContainerImpl) GetMedian() *int {
	if len(i.Slice) == 0 {
		return nil
	}

	sort.Ints(i.Slice)
	var index int
	if len(i.Slice)%2 == 0 { // its even
		index = len(i.Slice)/2 - 1
	} else { // its odd
		index = len(i.Slice) / 2
	}
	value := i.Slice[index]
	return &value
}
