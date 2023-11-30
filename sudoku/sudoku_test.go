package sudoku_test

import (
	"math"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ScarletTanager/godoku/sudoku"
)

var _ = Describe("Sudoku", func() {
	var (
		s        *sudoku.Sudoku
		initRows [][]int
		size     int
	)

	Describe("New", func() {
		It("Returns a new correctly-sized Sudoku", func() {
			for size := 2; size < 5; size++ {
				sud := sudoku.New(size)
				Expect(sud.SubsetSize).To(Equal(int(math.Pow(float64(size), 2))))
			}
		})
	})

	BeforeEach(func() {
		size = 3
		initRows = [][]int{
			{2, 7, 0, 4, 0, 8, 0, 0, 6},
			{4, 0, 0, 0, 6, 0, 0, 0, 3},
			{0, 9, 6, 7, 0, 2, 4, 0, 5},
			{0, 6, 4, 5, 0, 0, 0, 0, 0},
			{0, 0, 0, 9, 0, 0, 0, 0, 4},
			{7, 1, 0, 6, 4, 3, 8, 5, 2},
			{0, 0, 0, 0, 0, 5, 6, 0, 0},
			{5, 0, 0, 0, 0, 6, 0, 3, 9},
			{0, 2, 1, 0, 0, 4, 0, 0, 8},
		}
	})

	JustBeforeEach(func() {
		s = sudoku.NewFromRows(size, initRows)
		Expect(s).NotTo(BeNil())
	})

	Describe("Row", func() {
		var (
			j int
		)

		BeforeEach(func() {
			j = 5
		})

		It("Returns the row at the specified index", func() {
			Expect(s.Row(j)).To(ConsistOf(initRows[j]))
		})
	})

	Describe("Column", func() {
		var (
			k int
		)

		BeforeEach(func() {
			k = 3
		})

		It("Returns the column at the specified index", func() {
			for j, val := range s.Column(k) {
				Expect(val).To(Equal(initRows[j][k]))
			}
		})
	})

	Describe("Grid", func() {
		var (
			l int
		)

		BeforeEach(func() {
			l = 4
		})

		It("Returns the lth grid left to right, top to bottom", func() {
			Expect(s.Grid(l)).To(Equal(sudoku.Subset{5, 0, 0, 9, 0, 0, 6, 4, 3}))
		})
	})
})
