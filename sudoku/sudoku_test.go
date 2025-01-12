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

	rowsWithBlanks := [][]int{
		{0, 0, 7, 8, 0, 3, 0, 2, 9},
		{8, 9, 2, 0, 0, 5, 3, 0, 0},
		{1, 0, 3, 0, 7, 9, 0, 8, 0},
		{0, 8, 0, 6, 5, 0, 2, 0, 3},
		{7, 2, 0, 0, 0, 0, 9, 0, 0},
		{0, 0, 0, 0, 2, 0, 8, 0, 6},
		{0, 0, 8, 4, 0, 2, 6, 1, 0},
		{2, 0, 0, 0, 1, 0, 0, 9, 8},
		{0, 5, 1, 0, 0, 0, 0, 0, 2},
	}

	solvedRows := [][]int{
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

		It("Assigns the correct column indices", func() {
			s = sudoku.New(3)
			for j := 0; j < s.SubsetSize; j++ {
				for k := 0; k < s.SubsetSize; k++ {
					Expect(s.Row(j)[k].ColumnIndex).To(Equal(k))
				}
			}
		})
	})

	Describe("NewFromRows", func() {
		var (
			ss sudoku.Subset
		)

		BeforeEach(func() {
			size = 3
		})

		JustBeforeEach(func() {
			Expect(s).NotTo(BeNil())
		})

		When("Some squares are unset", func() {
			BeforeEach(func() {
				s = sudoku.NewFromRows(size, rowsWithBlanks)
			})

			// Test that the unset squares have any values set within the same row/column/subgrid
			// removed from their list of candidates/possible values
			It("Constrains the unset squares in each row", func() {
				for idx := 0; idx < s.SubsetSize; idx++ {
					ss = s.Row(idx)

					setValues := make([]int, 0)
					for _, gs := range ss {
						if gs.Value != 0 {
							setValues = append(setValues, gs.Value)
						}
					}

					for _, gs := range ss {
						if gs.Value == 0 {
							Expect(gs.Candidates()).NotTo(ContainElements(setValues))
						}
					}
				}
			})

			It("Constrains the unset squares in each column", func() {
				for idx := 0; idx < s.SubsetSize; idx++ {
					ss = s.Column(idx)

					setValues := make([]int, 0)
					for _, gs := range ss {
						if gs.Value != 0 {
							setValues = append(setValues, gs.Value)
						}
					}

					for _, gs := range ss {
						if gs.Value == 0 {
							Expect(gs.Candidates()).NotTo(ContainElements(setValues))
						}
					}
				}
			})

			It("Constrains the unset squares in each subgrid", func() {
				for idx := 0; idx < s.SubsetSize; idx++ {
					ss = s.Subgrid(idx)

					setValues := make([]int, 0)
					for _, gs := range ss {
						if gs.Value != 0 {
							setValues = append(setValues, gs.Value)
						}
					}

					for _, gs := range ss {
						if gs.Value == 0 {
							Expect(gs.Candidates()).NotTo(ContainElements(setValues))
						}
					}
				}
			})
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

	Describe("Subset methods", func() {
		var (
			ss sudoku.Subset
		)

		BeforeEach(func() {
			size = 3
			s = sudoku.NewFromRows(size, rowsWithBlanks)
		})

		JustBeforeEach(func() {
			Expect(s).NotTo(BeNil())
		})

		Describe("MaskValue", func() {
			var (
				val, squareIdx int
			)

			BeforeEach(func() {
				val = 6
				squareIdx = 7
				ss = s.Row(1)
				Expect(ss[4].Candidates()).To(ContainElement(val))
				Expect(ss[squareIdx].Candidates()).To(ContainElement(val))
			})

			It("Removes the value from the candidates of all squares other than the one specified", func() {
				ss.MaskValue(7, val)
				Expect(ss[4].Candidates()).NotTo(ContainElement(val))
				Expect(ss[squareIdx].Candidates()).To(ContainElement(val))
			})
		})

		Describe("AllValuesUnique", func() {
			BeforeEach(func() {
				s = sudoku.NewFromRows(size, solvedRows)
			})

			JustBeforeEach(func() {
				ss = s.Row(1)
			})

			When("The subset is composed of unique values", func() {
				It("Returns true", func() {
					Expect(ss.AllValuesUnique()).To(BeTrue())
				})
			})

			When("The subset contains duplicate values", func() {
				When("But the duplicate values are all zeroes", func() {
					BeforeEach(func() {
						s = sudoku.NewFromRows(size, rowsWithBlanks)
					})

					It("Returns true", func() {
						Expect(ss.AllValuesUnique()).To(BeTrue())
					})
				})

				When("And the duplicate values are nonzero", func() {
					JustBeforeEach(func() {
						ss[0].Value = 5
					})

					It("Returns false", func() {
						Expect(ss.AllValuesUnique()).To(BeFalse())
					})
				})
			})
		})
	})

	Describe("Sudoku methods", func() {
		BeforeEach(func() {
			size = 3
			initRows = rowsWithBlanks
		})

		JustBeforeEach(func() {
			s = sudoku.NewFromRows(size, initRows)
			Expect(s).NotTo(BeNil())
		})

		Describe("Solved", func() {
			When("Some squares do not have values", func() {
				It("Returns false", func() {
					Expect(s.Solved()).To(BeFalse())
				})
			})

			When("All squares have values", func() {
				BeforeEach(func() {
					// This is a valid solved sudoku from the NYT
					initRows = [][]int{
						{7, 8, 9, 2, 6, 4, 5, 1, 3},
						{2, 3, 6, 9, 1, 5, 7, 8, 4},
						{5, 4, 1, 8, 7, 3, 6, 2, 9},
						{8, 5, 7, 4, 9, 2, 3, 6, 1},
						{4, 6, 2, 7, 3, 1, 9, 5, 8},
						{1, 9, 3, 6, 5, 8, 2, 4, 7},
						{3, 2, 8, 5, 4, 7, 1, 9, 6},
						{9, 7, 4, 1, 2, 6, 8, 3, 5},
						{6, 1, 5, 3, 8, 9, 4, 7, 2},
					}
				})

				When("And no constraints are violated", func() {
					It("Returns true", func() {
						Expect(s.Solved()).To(BeTrue())
					})
				})

				When("The same value occurs twice within a column", func() {
					BeforeEach(func() {
						initRows[8][7] = 2
						initRows[8][8] = 7
					})

					It("Returns false", func() {
						Expect(s.Solved()).To(BeFalse())
					})
				})

				When("The same value occurs twice within a row", func() {
					BeforeEach(func() {
						initRows[7][8] = 2
						initRows[8][8] = 5
					})

					It("Returns false", func() {
						Expect(s.Solved()).To(BeFalse())
					})
				})
			})
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
					6, 5, 0,
					0, 0, 0,
					0, 2, 0}
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
					initRows = solvedRows
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
				j = 4
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
