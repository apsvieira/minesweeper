package views

import (
	"html/template"
	"net/http"

	"github.com/apsvieira/minesweeper/internal/game"
)

type GameView struct {
	templ *template.Template
}

func NewGameView(templ *template.Template) *GameView {
	return &GameView{templ: templ}
}

func (t *GameView) Render(w http.ResponseWriter, f *game.Field) {
	if err := t.templ.ExecuteTemplate(w, "game", f); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *GameView) RenderGrid(w http.ResponseWriter, f *game.Field) {
	if err := t.templ.ExecuteTemplate(w, "grid", f); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "text/html")
}
