package game

import (
	"fmt"
)

type Cell struct {
	// The number of mines adjacent to this cell
	numAdjacentMines int
	// Whether this cell is a mine
	isMine bool
	// Whether this cell has been revealed
	isRevealed bool
	// Whether this cell has been flagged
	isFlagged bool
}

// Make cells printable
func (c *Cell) String() string {
	if c.isFlagged {
		return "F"
	}
	if !c.isRevealed {
		return "."
	}
	if c.isMine {
		return "M"
	}
	return fmt.Sprintf("%d", c.numAdjacentMines)
}

func (c *Cell) reveal() {
	if c.isFlagged {
		return
	}

	c.isRevealed = true
}

func (c *Cell) revealAndCheck() error {
	c.reveal()
	if c.isMine {
		return ErrMineHit
	}

	return nil
}

func (c *Cell) flag() {
	c.isFlagged = true
}

func (c *Cell) unflag() {
	c.isFlagged = false
}
