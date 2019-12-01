package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

const maxN = 13

// CellColor is a cell color.
type CellColor byte

// CellLine is a line of CellColor.
type CellLine []CellColor

// CellMatrix is a matrix of CellColor.
type CellMatrix []CellLine

const (
	greyCell CellColor = iota
	blackCell
	whiteCell
)

// Kakurasu is a kakurasu.
type Kakurasu struct {
	n       int
	rowSums []int
	colSums []int
	m       CellMatrix
}

// return a CellMatrix with grey cells
func newCellMatrix(n int) CellMatrix {
	m := make([]CellLine, n)
	for i := range m {
		m[i] = make([]CellColor, n)
	}
	return m
}

// NewKakurasu return a new random Kakurasu.
func NewKakurasu(n int) *Kakurasu {
	if n > maxN {
		panic(fmt.Sprint("n too big:", n))
	}
	k := &Kakurasu{n: n, rowSums: make([]int, n), colSums: make([]int, n)}
	k.m = newCellMatrix(n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if rand.Intn(2) == 1 {
				k.rowSums[i] += j+1 // lines
				k.colSums[j] += i+1 // culumns
				k.m[i][j] = blackCell
			} else {
				k.m[i][j] = whiteCell
			}
		}
	}
	return k
}

// Print outputs the Kakurasu to stdOut.
func (k *Kakurasu) Print() {
	fmt.Print("  ")
	for j := 0; j < k.n; j++ {
		fmt.Printf("%2d", k.colSums[j])
	}
	for i := 0; i < k.n; i++ {
		fmt.Print("\n  ")
		for j := 0; j < k.n; j++ {
			switch k.m[i][j] {
			case blackCell:
				fmt.Print("**")
			case whiteCell:
				fmt.Print("  ")
			case greyCell:
				fmt.Print("??")
			default:
				panic("invalid cell color")
			}
		}
		fmt.Printf("\n%2d", k.rowSums[i])
		for j := 0; j < k.n; j++ {
			switch k.m[i][j] {
			case blackCell:
				fmt.Print("**")
			case whiteCell:
				fmt.Print("  ")
			case greyCell:
				fmt.Print("??")
			default:
				panic("invalid cell color")
			}
		}
	}
	fmt.Println()
}

func buildSolutions(n int) map[int][]CellLine {
	if n > maxN {
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

func solveLine(sols []CellLine) CellLine {
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
    		if sols[j][i] == blackCell{
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

// eliminate solutions incompatible with color c at index k
func filterLineSols(sols []CellLine, k int, c CellColor) []CellLine {
    var j int
	for i := range sols {
		if sols[i][k] == c {
			sols[j] = sols[i]
			j++
		}
	}
	return sols[:j]
}

func lineStr(l CellLine) string {
	c := make([]byte, len(l))
	for i := range l {
		switch l[i]{
		case blackCell:
		    c[i] = '*'
		case whiteCell:
		    c[i] = ' '
		default:
		    c[i] = '?'
		}
	}
	return string(c)
}


// Solve Kakurasu and return solution matrix.
func Solve(rowSums, colSums []int) CellMatrix {
	n := len(rowSums)
	if len(colSums) != n {
		panic("solve only square Kakurasu")
	}

	solTbl := buildSolutions(n)
	// printSolutions(solTbl)

	m := newCellMatrix(n)

	// copy solutions into rowSols and colSols
	rowSols := make([][]CellLine, n)
	colSols := make([][]CellLine, n)
	for i := 0; i < n; i++ {
		rowSols[i] = append(rowSols[i], solTbl[rowSums[i]]...)
		colSols[i] = append(colSols[i], solTbl[colSums[i]]...)
	}

	hasMod := true
	// as long as we have modified rows or columns
	for hasMod {
		hasMod = false
		// for rows
		for r := 0; r < n; r++ {
			sol := solveLine(rowSols[r])
			for c := 0; c < n; c++ {
				if sol[c] == greyCell {
					continue
				}
				if m[r][c] == greyCell {
					// we solved a new cell
					hasMod = true
					m[r][c] = sol[c]
					colSols[c] = filterLineSols(colSols[c], r, sol[c])
				} else {
					if m[r][c] != sol[c] {
						panic("matrix and solution mismatch")
					}
				}
			}
		}
		// for columns
		for c := 0; c < n; c++ {
			sol := solveLine(colSols[c])
			for r := 0; r < n; r++ {
				if sol[r] == greyCell {
					continue
				}
				if m[r][c] == greyCell {
					// we solved a new cell
					hasMod = true
					m[r][c] = sol[r]
					rowSols[r] = filterLineSols(rowSols[r], c, sol[r])
				} else {
					if m[r][c] != sol[r] {
						panic("matrix and solution mismatch")
					}
				}
			}
		}
	}
	return m
}

func main() {
	n := 5
	rand.Seed(time.Now().Unix())
	k := NewKakurasu(n)
	k.Print()

	fmt.Println("solution:")
	k2 := k
	k2.m = Solve(k.rowSums, k.colSums)
	k2.Print()

}
