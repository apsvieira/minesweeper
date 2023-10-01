package main

import (
	"fmt"

	"github.com/apsvieira/minesweeper/internal/game"
)

func main() {
	// Start a new game by asking the user for the dimensions and number of mines
	w, h, mines := getDimensions()
	f, err := game.NewField(w, h, mines)
	if err != nil {
		panic(err)
	}

	// The main game loop is just asking for coordinates and attempting to
	// reveal the cell. If the cell is a mine, the game is over.
	for !f.IsWon() {
		var x, y int
		fmt.Println("---")
		fmt.Println(f)
		fmt.Println("Enter x:")
		fmt.Scan(&x)
		fmt.Println("Enter y:")
		fmt.Scan(&y)
		if err := f.Reveal(x, y); err != nil {
			fmt.Println("You Lost!")
			f.RevealAll()
			fmt.Println(f)
			return
		}
	}

	fmt.Println("You Won!")
}

func getDimensions() (int, int, int) {
	fmt.Println("Enter width:")
	var width int
	fmt.Scan(&width)
	fmt.Println("Enter height:")
	var height int
	fmt.Scan(&height)
	fmt.Println("Enter number of mines:")
	var numMines int
	fmt.Scan(&numMines)

	return width, height, numMines
}
