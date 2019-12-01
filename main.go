package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

const maxLineLength = 13

// CellColor is a cell color.
type CellColor byte

const (
	greyCell CellColor = iota
	blackCell
	whiteCell
)

// CellLine is a line of CellColor.
type CellLine []CellColor

// CellGrid is an n rows by m columns grid of CellColor.
type CellGrid struct {
	n, m int
	g    []CellLine
}

// NewCellGrid returns a CellMatrix with grey cells.
func NewCellGrid(n, m int) *CellGrid {
	g := make([]CellLine, n)
	for i := range g {
		g[i] = make(CellLine, m)
	}
	return &CellGrid{n: n, m: m, g: g}
}

// Clone return a copy of the CellGrid.
func (g *CellGrid) Clone() *CellGrid {
	newG := &CellGrid{n: g.n, m: g.m, g: make([]CellLine, g.n)}
	for i := range newG.g {
		newG.g[i] = append(newG.g[i], g.g[i]...)
	}
	return newG
}

// FillRandomly fills grid randomly with black and white cell colors.
func (g *CellGrid) FillRandomly() {
	for i := 0; i < g.n; i++ {
		for j := 0; j < g.m; j++ {
			if rand.Intn(2) == 1 {
				g.g[i][j] = blackCell
			} else {
				g.g[i][j] = whiteCell
			}
		}
	}
}

// ComputeSums of black cell weight in rows and columns.
func (g *CellGrid) ComputeSums() (rowSums, colSums []int) {
	rowSums = make([]int, g.n)
	colSums = make([]int, g.m)
	for i := 0; i < g.n; i++ {
		for j := 0; j < g.m; j++ {
			if g.g[i][j] == blackCell {
				rowSums[i] += j + 1
				colSums[j] += i + 1
			}
		}
	}
	return
}

// Print outputs CellGrid to stdOut.
func (g *CellGrid) Print() {
	for i := 0; i < g.n; i++ {
		for j := 0; j < g.m; j++ {
			switch g.g[i][j] {
			case blackCell:
				fmt.Print("*")
			case whiteCell:
				fmt.Print(" ")
			case greyCell:
				fmt.Print("?")
			default:
				panic("invalid cell color")
			}
		}
		fmt.Println()
	}
}

//PrintSums prints the sums to stdOut.
func PrintSums(rowSums, colSums []int) {
	fmt.Print("row sums:")
	for i := range rowSums {
		fmt.Printf(" %2d", rowSums[i])
	}
	fmt.Println()
	fmt.Print("col sums:")
	for i := range colSums {
		fmt.Printf(" %2d", colSums[i])
	}
	fmt.Println()
}

func buildSolutions(n int) map[int][]CellLine {
	if n > maxLineLength {
		panic(fmt.Sprint("n too big:", n))
	}
	m := make(map[int][]CellLine)
	nSol := 1 << n
	for i := 0; i < nSol; i++ {
		var sum int
		v := i
		sol := make(CellLine, n)
		for k := 0; k < n; k++ {
			if v&1 != 0 {
				sum += k + 1
				sol[k] = blackCell
			} else {
				sol[k] = whiteCell
			}
			v >>= 1
		}
		m[sum] = append(m[sum], sol)
	}
	return m
}

func printSolutions(solTbl map[int][]CellLine) {
	sums := make([]int, 0, len(solTbl))
	for k := range solTbl {
		sums = append(sums, k)
	}
	sort.IntSlice(sums).Sort()
	for _, k := range sums {
		v := solTbl[k]
		fmt.Println("sum:", k)
		for i, l := range v {
			fmt.Printf("%2d. ", i)
			for j := range l {
				switch l[j] {
				case greyCell:
					fmt.Print("?")
				case blackCell:
					fmt.Print("*")
				case whiteCell:
					fmt.Print(" ")
				default:
					panic("invalid cell color")
				}
			}
			fmt.Println()
		}
	}
}

func deduceColorsFormSols(sols []CellLine) CellLine {
	nSols := len(sols)
	if nSols == 0 {
		panic("line has no solutions")
	}
	if nSols == 1 {
		return sols[0]
	}
	n := len(sols[0])
	c := make(CellLine, n)
	for i := range c {
		var acc int
		for j := range sols {
			if sols[j][i] == blackCell {
				acc++
			}
		}
		switch acc {
		case 0:
			c[i] = whiteCell
		case nSols:
			c[i] = blackCell
		default:
			c[i] = greyCell
		}
	}
	return c
}

// eliminate solutions incompatible with black or white color c at index k.
func pruneSols(sols []CellLine, k int, c CellColor) []CellLine {
	var j int
	for i := range sols {
		if sols[i][k] == c {
			sols[j] = sols[i]
			j++
		}
	}
	return sols[:j]
}

func initSols(sums []int, m int) [][]CellLine {
	n := len(sums)
	solTbl := buildSolutions(m)
	sols := make([][]CellLine, n)
	for i := 0; i < n; i++ {
		sols[i] = append(sols[i], solTbl[sums[i]]...)
	}
	return sols
}

func copySols(sols [][]CellLine) [][]CellLine {
	newSols := make([][]CellLine, len(sols))
	for i := range newSols {
		newSols[i] = append(newSols[i], sols[i]...)
	}
	return newSols
}

type solveState struct {
	rowSols [][]CellLine
	colSols [][]CellLine
	g       *CellGrid
	nGreys  int
}

func newSolveState(rowSums, colSums []int) *solveState {
	return &solveState{
		rowSols: initSols(rowSums, len(colSums)),
		colSols: initSols(colSums, len(rowSums)),
		g:       NewCellGrid(len(rowSums), len(colSums)),
		nGreys:  len(rowSums) * len(colSums),
	}
}

func (s *solveState) clone() *solveState {
	return &solveState{
		rowSols: copySols(s.rowSols),
		colSols: copySols(s.colSols),
		g:       s.g.Clone(),
		nGreys:  s.nGreys,
	}
}

func (s *solveState) findFirstGrey() (r, c int) {
	for r := 0; r < s.g.n; r++ {
		for c := 0; c < s.g.m; c++ {
			if s.g.g[r][c] == greyCell {
				return r, c
			}
		}
	}
	panic("no grey cell found")
}

func (s *solveState) setFirstGreyColorTo(clr CellColor) {
	r, c := s.findFirstGrey()
	s.g.g[r][c] = clr
	s.nGreys--
	s.rowSols[r] = pruneSols(s.rowSols[r], c, clr)
	s.colSols[c] = pruneSols(s.colSols[c], r, clr)
}

func (s *solveState) deduce() {
	hasMod := true
	// as long as we have modified rows or columns
	for hasMod && s.nGreys > 0 {
		hasMod = false
		// for each row
		for r := 0; r < s.g.n; r++ {
			sol := deduceColorsFormSols(s.rowSols[r])
			for c := 0; c < s.g.m; c++ {
				if sol[c] == greyCell {
					continue
				}
				if s.g.g[r][c] == greyCell {
					// we solved a new cell
					hasMod = true
					s.nGreys--
					s.g.g[r][c] = sol[c]
					s.colSols[c] = pruneSols(s.colSols[c], r, sol[c])
				} else {
					if s.g.g[r][c] != sol[c] {
						panic("matrix and solution mismatch")
					}
				}
			}
		}
		// for each column
		for c := 0; c < s.g.m; c++ {
			sol := deduceColorsFormSols(s.colSols[c])
			for r := 0; r < s.g.n; r++ {
				if sol[r] == greyCell {
					continue
				}
				if s.g.g[r][c] == greyCell {
					// we solved a new cell
					hasMod = true
					s.nGreys--
					s.g.g[r][c] = sol[r]
					s.rowSols[r] = pruneSols(s.rowSols[r], c, sol[r])
				} else {
					if s.g.g[r][c] != sol[r] {
						panic("matrix and solution mismatch")
					}
				}
			}
		}
	}
}

func (s *solveState) doSolve(sols []*CellGrid) []*CellGrid {
	s.deduce()
	if s.nGreys == 0 {
		return append(sols, s.g)
	}
	s2 := s.clone()
	s.setFirstGreyColorTo(blackCell)
	s2.setFirstGreyColorTo(whiteCell)
	return s2.doSolve(s.doSolve(sols))
}

// Solve Kakurasu and return list of solution grids.
func Solve(rowSums, colSums []int) []*CellGrid {
	s := newSolveState(rowSums, colSums)
	return s.doSolve([]*CellGrid{})
}

func main() {
	rand.Seed(time.Now().Unix())
	g := NewCellGrid(5, 6)
	g.FillRandomly()
	rowSums, colSums := g.ComputeSums()

	fmt.Println("Kakurasu:")
	g.Print()
	PrintSums(rowSums, colSums)
	fmt.Println()

	sols := Solve(rowSums, colSums)
	fmt.Println("Solutions:")
	for i := range sols {
		sols[i].Print()
		fmt.Println()
	}
}
