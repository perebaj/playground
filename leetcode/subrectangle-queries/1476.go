package main

import "fmt"

type SubrectangleQueries struct {
	matrix *[][]int
}

func Constructor(rectangle [][]int) SubrectangleQueries {
	return SubrectangleQueries{
		matrix: &rectangle,
	}

}

func (s *SubrectangleQueries) NumRows() int {
	return len(*s.matrix)
}

func (s *SubrectangleQueries) NumColums() int {
	return len((*s.matrix)[0])
}

func (s *SubrectangleQueries) UpdateSubrectangle(row1 int, col1 int, row2 int, col2 int, newValue int) {
	maxNumRows := s.NumRows()
	maxNumCols := s.NumColums()
	if row2 >= maxNumRows {
		row2 = maxNumCols - 1
	}

	if col2 >= maxNumCols {
		row2 = maxNumCols - 1
	}

	for irow := row1; irow <= row2; irow++ {
		for jcol := col1; jcol <= col2; jcol++ {
			(*s.matrix)[irow][jcol] = newValue
		}
	}
}

func (s *SubrectangleQueries) GetValue(row int, col int) int {
	if row > s.NumRows() || col > s.NumColums() {
		return -1
	} else {
		return (*s.matrix)[row][col]
	}
}

func main() {
	n := Constructor([][]int{
		{1, 2, 1},
		{4, 3, 4},
		{3, 2, 1},
	})

	n.UpdateSubrectangle(0, 0, 2, 2, 100)
	fmt.Println((*n.matrix))
}
