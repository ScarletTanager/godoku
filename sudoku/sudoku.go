package sudoku

import (
	"math"
)

type Sudoku struct {
	SubsetSize int
	sizeFactor int
	values     []int
}

type Subset []int

func New(size int) *Sudoku {
	subsetSize := int(math.Pow(float64(size), 2))

	return &Sudoku{
		SubsetSize: subsetSize,
		sizeFactor: size,
		values:     make([]int, subsetSize*subsetSize),
	}
}

func NewFromRows(size int, rows [][]int) *Sudoku {
	s := New(size)
	if len(rows) != s.SubsetSize {
		return nil
	}

	for j, row := range rows {
		if len(row) != s.SubsetSize {
			return nil
		}

		for k, val := range row {
			s.values[(j*s.SubsetSize)+k] = val
		}
	}

	return s
}

func (s *Sudoku) Set(j, k, value int) {
	s.values[(j*s.SubsetSize)+k] = value
}

func (s *Sudoku) Row(j int) Subset {
	return s.values[j*s.SubsetSize : (j+1)*s.SubsetSize]
}

func (s *Sudoku) Column(k int) Subset {
	col := make([]int, s.SubsetSize)
	for j := 0; j < s.SubsetSize; j++ {
		col[j] = s.values[k+(j*s.SubsetSize)]
	}

	return col
}

func (s *Sudoku) Grid(l int) Subset {
	grid := make([]int, s.SubsetSize)
	offset := l % s.sizeFactor
	startIndex := ((l - offset) * s.SubsetSize) + (offset * s.sizeFactor)
	for gri := 0; gri < s.sizeFactor; gri++ {
		for gci := 0; gci < s.sizeFactor; gci++ {
			grid[(gri*s.sizeFactor)+gci] = s.values[startIndex+gci+(gri*s.SubsetSize)]
		}
	}

	return grid
}
