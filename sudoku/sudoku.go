package sudoku

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
)

type Sudoku struct {
	SubsetSize int
	sizeFactor int
	Values     []*GridSquare
}

// New returns a Sudoku with all squares set to 0, any square
// can assume any value.  The overall Sudoku is (size^2) x (size^2) -
// in other words, it is a sizexsize grid of sizexsize subgrids.
func New(size int) *Sudoku {
	subsetSize := int(math.Pow(float64(size), 2))

	s := &Sudoku{
		SubsetSize: subsetSize,
		sizeFactor: size,
		Values:     make([]*GridSquare, subsetSize*subsetSize),
	}

	for i := range s.Values {
		s.Values[i] = &GridSquare{}

		// Initialize the row, column, and subgrid indices
		j := i / subsetSize
		k := i % subsetSize
		s.Values[i].RowIndex = j
		s.Values[i].ColumnIndex = k
		s.Values[i].SubgridIndex = SubgridIndex(j, k, size)

		// Initialize each grid square so that it could potentially contain any value
		s.Values[i].PossibleValues = make(map[int]struct{})
		for pv := 1; pv <= subsetSize; pv++ {
			s.Values[i].PossibleValues[pv] = struct{}{}
		}
	}

	return s
}

// NewFromRows returns a Sudoku with the Values set, any unset
// square has its candidates constrained.
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
			s.Values[(j*s.SubsetSize)+k].Value = val
		}
	}

	return s
}

// Solved returns true if we have solved the puzzle, false otherwise.
func (s *Sudoku) Solved() bool {
	// Does every square have a legal value?
	for _, gridSquare := range s.Values {
		if gridSquare.Value < 1 || gridSquare.Value > 9 {
			return false
		}
	}

	return s.obeysConstraints()
}

func (s *Sudoku) obeysConstraints() bool {
	for i := 0; i < s.SubsetSize; i++ {
		if !s.Row(i).allValuesUnique() {
			return false
		}

		if !s.Column(i).allValuesUnique() {
			return false
		}

		if !s.Subgrid(i).allValuesUnique() {
			return false
		}
	}

	return true
}

// Row returns a Subset containing the squares contained by the
// column at position k, indexed from 0 at the top to (size^2)-1 at the bottom.
func (s *Sudoku) Row(j int) Subset {
	return s.Values[j*s.SubsetSize : (j+1)*s.SubsetSize]
}

// Column returns a Subset containing the squares contained by the
// column at position k, indexed from 0 at the left to (size^2)-1 at the right.
func (s *Sudoku) Column(k int) Subset {
	col := make(Subset, s.SubsetSize)
	for j := 0; j < s.SubsetSize; j++ {
		col[j] = s.Values[k+(j*s.SubsetSize)]
	}

	return col
}

// Subgrid returns a Subset containing the grid squares contained by the
// subgrid at position l, with 0 at the top left, (size^2)-1 at the bottom
// right, left to right, top to bottom
func (s *Sudoku) Subgrid(l int) Subset {
	grid := make(Subset, s.SubsetSize)
	offset := l % s.sizeFactor
	startIndex := ((l - offset) * s.SubsetSize) + (offset * s.sizeFactor)
	for gri := 0; gri < s.sizeFactor; gri++ {
		for gci := 0; gci < s.sizeFactor; gci++ {
			grid[(gri*s.sizeFactor)+gci] = s.Values[startIndex+gci+(gri*s.SubsetSize)]
		}
	}

	return grid
}

// Set assigns the Value of the GridSquare at row j, column k to val.
// It then removes val from the PossibleValues for every other square
// within the same row, column, or gridsquare.
func (s *Sudoku) Set(j, k, val int) error {
	gridSquare := s.Values[(j*s.SubsetSize)+k]

	if _, ok := gridSquare.PossibleValues[val]; !ok {
		pvs := make([]int, 0)
		for v, _ := range gridSquare.PossibleValues {
			pvs = append(pvs, v)
		}

		return fmt.Errorf("Value %d not in list of candidates %v", val, pvs)
	}

	// Set the value and set the possibles to only that value
	gridSquare.Value = val
	gridSquare.PossibleValues = map[int]struct{}{val: struct{}{}}

	// Remove the value from the candidates for every other square in the same
	// row, column, or subgrid
	s.Row(gridSquare.RowIndex).MaskValue(k, val)
	s.Column(gridSquare.ColumnIndex).MaskValue(j, val)
	s.Subgrid(gridSquare.SubgridIndex).MaskValue(SquareIndexInSubgrid(j, k, s.sizeFactor), val)

	return nil
}

// String prints the string representation of the sudoku
func (s *Sudoku) String() string {
	return ""
}

// PrintCurrent prints the current state of the sudoku
func (s *Sudoku) Current() string {
	var current strings.Builder

	for i := 0; i < s.SubsetSize; i++ {
		r := s.Row(i)
		for _, gs := range r {
			current.WriteString(gs.String() + " ")
		}
		current.WriteString("\n")
	}

	return current.String()
}

// Hash returns the sha256 checksum of the sudoku's current state
func (s *Sudoku) Hash() [sha256.Size]byte {
	b, _ := json.Marshal(s)
	return sha256.Sum256(b)
}

// SubgridIndex returns the index of the subgrid containing the square at
// row j, column k, given a sizexsize grid of sizexsize subgrids, with the
// subgrids indexed starting with 0 at the top left, left to right, top to bottom.
func SubgridIndex(j, k, size int) int {
	return j - (j % size) + (k / size)
}

// SquareIndexInSubgrid returns the index of the square _within the subgrid_ -
// indexed from 0 at top left to (size^2)-1 at the bottom right, left to right,
// top to bottom.
func SquareIndexInSubgrid(j, k, size int) int {
	return ((j % size) * size) + (k % size)
}

type Possibles map[int]struct{}

func (p Possibles) MarshalJSON() ([]byte, error) {
	pvs := make([]int, 0)
	for pv, _ := range p {
		pvs = append(pvs, pv)
	}
	sort.Ints(pvs)
	return json.Marshal(pvs)
}

func (p *Possibles) UnmarshalJSON(b []byte) error {
	var (
		pvs []int
	)

	if err := json.Unmarshal(b, &pvs); err != nil {
		return err
	}

	for _, pv := range pvs {
		(*p)[pv] = struct{}{}
	}

	return nil
}

// GridSquare is an individual square within the overall grid.  Each
// Sudoku will contain n^4 squares (where the Sudoku is composed of n^2 nxn subgrids).
type GridSquare struct {
	Value          int
	PossibleValues Possibles
	RowIndex       int
	ColumnIndex    int
	SubgridIndex   int
}

// Constrain limits the candidate Values for the square to those contained
// by possibles
func (gs *GridSquare) Constrain(possibles []int) {
	candidates := make(map[int]struct{})
	for _, pv := range possibles {
		candidates[pv] = struct{}{}
		// Go ahead and make sure it's in the stored list of possibles
		gs.PossibleValues[pv] = struct{}{}
	}

	for v, _ := range gs.PossibleValues {
		// If the value is not in the list of candidates passed in, delete it
		if _, ok := candidates[v]; !ok {
			delete(gs.PossibleValues, v)
		}
	}
}

// RemoveCandidate removes the value exclusion from PossibleValues
func (gs *GridSquare) RemoveCandidate(exclusion int) {
	delete(gs.PossibleValues, exclusion)
}

func (gs *GridSquare) String() string {
	if gs.hasLegalValue() {
		return fmt.Sprintf("**%d**", gs.Value)
	}

	pvList := make([]int, 0)
	for pv, _ := range gs.PossibleValues {
		pvList = append(pvList, pv)
	}

	return fmt.Sprintf("%v", pvList)
}

func (gs *GridSquare) Current() string {
	return ""
}

func (gs *GridSquare) hasLegalValue() bool {
	return gs.Value > 0 && gs.Value < 10
}

type Subset []*GridSquare

// MaskValue removes val from the PossibleValues of every GridSquare in the
// Subset other than the one at position i
func (ss Subset) MaskValue(i, val int) {
	for ssIndex := range ss {
		if ssIndex != i {
			// Delete the value from the candidate list
			ss[ssIndex].RemoveCandidate(val)
		}
	}
}

func (ss Subset) allValuesUnique() bool {
	ssVals := make(map[int]int)

	for _, i := range ss {
		if i.Value > 0 {
			if ssVals[i.Value] > 0 {
				return false
			} else {
				ssVals[i.Value] = 1
			}
		}
	}

	return true
}
