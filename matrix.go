package femtioelva

import (
	"fmt"
	"strings"
)

type Matrix struct {
	rows int
	cols int
	data []int // row major
}

func NewMatrix(rows, cols int) Matrix {
	return Matrix{
		rows: rows,
		cols: cols,
		data: make([]int, rows*cols),
	}
}

func (m Matrix) Width() int {
	return m.cols
}

func (m Matrix) Height() int {
	return m.rows
}

func (m Matrix) String() string {
	var s strings.Builder
	for i, v := range m.data {
		if i > 0 && i%m.cols == 0 {
			fmt.Fprintf(&s, "\n")
		}
		if v > 0 {
			fmt.Fprintf(&s, "%v", v)
		} else {
			fmt.Fprintf(&s, ".")
		}
	}
	return s.String()
}

func (m Matrix) At(row, col int) int {
	index := row*m.cols + col
	return m.data[index]
}

func (m Matrix) Incr(row, col int) {
	index := row*m.cols + col
	m.data[index]++
}

func (m Matrix) Set(row, col, val int) {
	index := row*m.cols + col
	m.data[index] = val
}
