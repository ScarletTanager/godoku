package sudoku_test

import (
	"math"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ScarletTanager/godoku/sudoku"
)

var _ = Describe("Sudoku", func() {
	var (
		s          *sudoku.Sudoku
		initRows   [][]int
		size       int
		subsetSize int
	)

	BeforeEach(func() {
		s = nil
	})

	JustBeforeEach(func() {
		subsetSize = int(math.Pow(float64(size), 2))
	})

	Describe("New", func() {
		It("Returns a new correctly-sized Sudoku", func() {
			for size := 2; size < 5; size++ {
				s = sudoku.New(size)
				Expect(s.SubsetSize).To(Equal(int(math.Pow(float64(size), 2))))
			}
		})

		It("Assigns the correct row indices", func() {
			s = sudoku.New(3)
			for j := 0; j < s.SubsetSize; j++ {
				for k := 0; k < s.SubsetSize; k++ {
					Expect(s.Row(j)[k].RowIndex).To(Equal(j))
				}
			}
		})
	})

	Describe("SubgridIndex", func() {
		var (
			subgridIndices [][]int
		)

		BeforeEach(func() {
			size = 3
			subgridIndices = [][]int{
				{0, 0, 0, 1, 1, 1, 2, 2, 2},
				{0, 0, 0, 1, 1, 1, 2, 2, 2},
				{0, 0, 0, 1, 1, 1, 2, 2, 2},
				{3, 3, 3, 4, 4, 4, 5, 5, 5},
				{3, 3, 3, 4, 4, 4, 5, 5, 5},
				{3, 3, 3, 4, 4, 4, 5, 5, 5},
				{6, 6, 6, 7, 7, 7, 8, 8, 8},
				{6, 6, 6, 7, 7, 7, 8, 8, 8},
				{6, 6, 6, 7, 7, 7, 8, 8, 8},
			}
		})

		It("Returns the correct subgrid index for the specified row and column", func() {
			for j := 0; j < subsetSize; j++ {
				for k := 0; k < subsetSize; k++ {
					Expect(sudoku.SubgridIndex(j, k, size)).To(Equal(subgridIndices[j][k]))
				}
			}
		})
	})

	Describe("SquareIndexInSubgrid", func() {
		var (
			squareIndices [][]int
		)

		BeforeEach(func() {
			size = 3
			squareIndices = [][]int{
				{0, 1, 2, 0, 1, 2, 0, 1, 2},
				{3, 4, 5, 3, 4, 5, 3, 4, 5},
				{6, 7, 8, 6, 7, 8, 6, 7, 8},
				{0, 1, 2, 0, 1, 2, 0, 1, 2},
				{3, 4, 5, 3, 4, 5, 3, 4, 5},
				{6, 7, 8, 6, 7, 8, 6, 7, 8},
				{0, 1, 2, 0, 1, 2, 0, 1, 2},
				{3, 4, 5, 3, 4, 5, 3, 4, 5},
				{6, 7, 8, 6, 7, 8, 6, 7, 8},
			}
		})

		It("Returns the correct index within the subgrid for the specified row and column", func() {
			for j := 0; j < subsetSize; j++ {
				for k := 0; k < subsetSize; k++ {
					Expect(sudoku.SquareIndexInSubgrid(j, k, size)).To(Equal(squareIndices[j][k]))
				}
			}
		})
	})

	Describe("Sudoku methods", func() {
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

		Describe("Solved", func() {

		})

		Describe("Row", func() {
			var (
				j int
			)

			BeforeEach(func() {
				j = 5
			})

			It("Returns the row at the specified index", func() {
				for k, gs := range s.Row(j) {
					Expect(gs.Value).To(Equal(initRows[j][k]))
				}
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
					Expect(val.Value).To(Equal(initRows[j][k]))
				}
			})
		})

		Describe("Subgrid", func() {
			var (
				l            int
				expectedVals []int
			)

			BeforeEach(func() {
				l = 4
				expectedVals = []int{
					5, 0, 0,
					9, 0, 0,
					6, 4, 3}
			})

			It("Returns the lth grid left to right, top to bottom", func() {
				for i, gs := range s.Subgrid(l) {
					Expect(gs.Value).To(Equal(expectedVals[i]))
				}
			})
		})

		Describe("Solved", func() {
			When("There are still 0-valued squares", func() {
				It("Returns false", func() {
					Expect(s.Solved()).To(BeFalse())
				})
			})

			When("There are duplicate values within a constrained subset (row/column/subgrid)", func() {
				BeforeEach(func() {
					initRows[2][0] = 2
				})

				It("Returns false", func() {
					Expect(s.Solved()).To(BeFalse())
				})
			})

			When("The puzzle has been solved", func() {
				BeforeEach(func() {
					initRows = [][]int{
						{1, 9, 4, 2, 6, 8, 3, 7, 5},
						{8, 5, 3, 9, 7, 4, 6, 1, 2},
						{7, 2, 6, 3, 5, 1, 9, 8, 4},
						{5, 1, 7, 8, 4, 3, 2, 9, 6},
						{9, 4, 8, 1, 2, 6, 7, 5, 3},
						{3, 6, 2, 7, 9, 5, 8, 4, 1},
						{4, 7, 5, 6, 3, 9, 1, 2, 8},
						{2, 3, 1, 5, 8, 7, 4, 6, 9},
						{6, 8, 9, 4, 1, 2, 5, 3, 7},
					}
				})

				It("Returns true", func() {
					Expect(s.Solved()).To(BeTrue())
				})
			})
		})

		Describe("Set", func() {
			var (
				j, k, val int
			)

			BeforeEach(func() {
				j = 3
				k = 4
				val = 8
			})

			JustBeforeEach(func() {
				Expect(s.Row(j)[k].Value).NotTo(Equal(val))
			})

			It("Sets the specified square to contain the specified value", func() {
				s.Set(j, k, val)
				Expect(s.Row(j)[k].Value).To(Equal(val))
			})

			It("Does not return an error", func() {
				Expect(s.Set(j, k, val)).NotTo(HaveOccurred())
			})

			When("The value is not in the list of candidates for the square", func() {
				JustBeforeEach(func() {
					s.Row(j)[k].RemoveCandidate(val)
				})

				It("Returns an error", func() {
					Expect(s.Set(j, k, val)).To(HaveOccurred())
				})

				It("Does not modify the value", func() {
					prev := s.Row(j)[k].Value
					s.Set(j, k, val)
					Expect(s.Row(j)[k].Value).To(Equal(prev))
					Expect(s.Row(j)[k].Value).NotTo(Equal(val))
				})
			})
		})
	})
})
