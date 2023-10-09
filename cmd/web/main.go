package main

import (
	"net/http"

	server "github.com/apsvieira/minesweeper/adapters/http"
)

func main() {
	// Create a new router and serve.
	h, err := server.NewHandler()
	if err != nil {
		panic(err)
	}

	http.ListenAndServe(":8080", h)
}
