package game

import (
	"fmt"
)

type Cell struct {
	// The number of mines adjacent to this cell
	NumAdjacentMines int
	// Whether this cell is a mine
	IsMine bool
	// Whether this cell has been revealed
	IsRevealed bool
	// Whether this cell has been flagged
	IsFlagged bool
}

// Make cells printable
func (c *Cell) String() string {
	if c.IsFlagged {
		return "F"
	}
	if !c.IsRevealed {
		return "."
	}
	if c.IsMine {
		return "M"
	}
	return fmt.Sprintf("%d", c.NumAdjacentMines)
}

func (c *Cell) reveal() {
	if c.IsFlagged {
		return
	}

	c.IsRevealed = true
}

func (c *Cell) revealAndCheck() error {
	c.reveal()
	if c.IsMine {
		return ErrMineHit
	}

	return nil
}

func (c *Cell) flag() {
	c.IsFlagged = true
}

func (c *Cell) unflag() {
	c.IsFlagged = false
}
