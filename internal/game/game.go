package game

import (
	"math/rand"
)

type Field struct {
	// The Cells is a 2D array of Cells
	Cells  [][]*Cell
	width  int
	height int

	// The number of mines in the field
	numMines int
	// The number of mines that have been flagged
	numFlaggedMines int
}

// Create a new field with the given dimensions and number of mines
func NewField(width, height, numMines int) (*Field, error) {
	if numMines < 0 {
		return nil, &ErrNegativeMines{numMines}
	}

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
		cells[y][x].IsMine = true
	}

	f := &Field{
		Cells:           cells,
		width:           width,
		height:          height,
		numMines:        numMines,
		numFlaggedMines: 0,
	}
	f.calculateAdjacentMines()

	return f, nil
}

// Make fields printable
func (f *Field) String() string {
	s := ""
	for _, row := range f.Cells {
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
		return &ErrOutOfBounds{x, y, f.width, f.height}
	}

	cell := f.Cells[y][x]
	return cell.revealAndCheck()
}

func (f *Field) isOutOfBounds(x, y int) bool {
	return x < 0 || x >= f.width || y < 0 || y >= f.height
}

func (f *Field) Flag(x, y int) {
	cell := f.Cells[y][x]
	cell.flag()
	if cell.IsMine {
		f.numFlaggedMines++
	}
}

func (f *Field) Unflag(x, y int) {
	cell := f.Cells[y][x]
	cell.unflag()
	if cell.IsMine {
		f.numFlaggedMines--
	}
}

func (f *Field) IsWon() bool {
	return f.numFlaggedMines == f.numMines
}

func (f *Field) RevealAll() {
	for _, row := range f.Cells {
		for _, cell := range row {
			cell.reveal()
		}
	}
}

func (f *Field) calculateAdjacentMines() {
	for y, row := range f.Cells {
		for x, cell := range row {
			if cell.IsMine {
				continue
			}
			for _, neighbor := range f.getNeighbors(x, y) {
				if neighbor.IsMine {
					cell.NumAdjacentMines++
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
			if nx < 0 || nx >= len(f.Cells[0]) || ny < 0 || ny >= len(f.Cells) {
				continue
			}
			neighbors = append(neighbors, f.Cells[ny][nx])
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
