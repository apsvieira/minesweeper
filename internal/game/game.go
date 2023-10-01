package game

import (
	"math/rand"
)

type Field struct {
	// The cells is a 2D array of cells
	cells  [][]*Cell
	width  int
	height int

	// The number of mines in the field
	numMines int
	// The number of mines that have been flagged
	numFlaggedMines int
}

// Create a new field with the given dimensions and number of mines
func NewField(width, height, numMines int) (*Field, error) {
	gridSize := width * height
	if numMines > gridSize {
		return nil, &ErrTooManyMines{numMines, gridSize}
	}

	cells := emptyCellGrid(height, width)
	mines := make([]int, numMines)
	for i := 0; i < numMines; i++ {
		mines[i] = rand.Intn(gridSize)
	}

	for _, mine := range mines {
		x := mine % width
		y := mine / width
		cells[y][x].isMine = true
	}

	f := &Field{cells, numMines, 0, width, height}
	f.calculateAdjacentMines()

	return f, nil
}

// Make fields printable
func (f *Field) String() string {
	s := ""
	for _, row := range f.cells {
		for _, cell := range row {
			s += cell.String()
		}
		s += "\n"
	}
	return s
}

// Attempt to reveal cell at position (x, y)
func (f *Field) Reveal(x, y int) error {
	if f.isOutOfBounds(x, y) {
		return &ErrOutOfBounds{x, y, len(f.cells[0]), len(f.cells)}
	}

	cell := f.cells[y][x]
	return cell.revealAndCheck()
}

func (f *Field) isOutOfBounds(x, y int) bool {
	return x < 0 || x >= f.width || y < 0 || y >= f.height
}

func (f *Field) Flag(x, y int) {
	cell := f.cells[y][x]
	cell.flag()
	if cell.isMine {
		f.numFlaggedMines++
	}
}

func (f *Field) Unflag(x, y int) {
	cell := f.cells[y][x]
	cell.unflag()
	if cell.isMine {
		f.numFlaggedMines--
	}
}

func (f *Field) IsWon() bool {
	return f.numFlaggedMines == f.numMines
}

func (f *Field) RevealAll() {
	for _, row := range f.cells {
		for _, cell := range row {
			cell.reveal()
		}
	}
}

func (f *Field) calculateAdjacentMines() {
	for y, row := range f.cells {
		for x, cell := range row {
			if cell.isMine {
				continue
			}
			for _, neighbor := range f.getNeighbors(x, y) {
				if neighbor.isMine {
					cell.numAdjacentMines++
				}
			}
		}
	}
}

func (f *Field) getNeighbors(x, y int) []*Cell {
	neighbors := make([]*Cell, 0, 8)
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			nx := x + i
			ny := y + j
			if nx < 0 || nx >= len(f.cells[0]) || ny < 0 || ny >= len(f.cells) {
				continue
			}
			neighbors = append(neighbors, f.cells[ny][nx])
		}
	}
	return neighbors
}

func emptyCellGrid(height int, width int) [][]*Cell {
	cells := make([][]*Cell, height)
	for i := 0; i < height; i++ {
		cells[i] = make([]*Cell, width)
		for j := 0; j < width; j++ {
			cells[i][j] = &Cell{}
		}
	}
	return cells
}
