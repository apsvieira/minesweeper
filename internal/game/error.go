package game

import (
	"errors"
	"fmt"
)

var ErrMineHit = errors.New("mine hit")

type ErrTooManyMines struct {
	numMines int
	gridSize int
}

func (e *ErrTooManyMines) Error() string {
	return fmt.Sprintf("too many mines (%d) for grid size (%d)", e.numMines, e.gridSize)
}

type ErrOutOfBounds struct {
	x    int
	y    int
	xMax int
	yMax int
}

func (e *ErrOutOfBounds) Error() string {
	return fmt.Sprintf("position (%d, %d) is out of bounds (max: %d, %d)", e.x, e.y, e.xMax, e.yMax)
}
