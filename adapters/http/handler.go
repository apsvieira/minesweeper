package http

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/apsvieira/minesweeper/adapters/http/views"
	"github.com/apsvieira/minesweeper/internal/game"
)

func NewHandler() (http.Handler, error) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.NoCache)
	r.Use(middleware.Heartbeat("/health"))

	staticFiles := http.FileServer(http.Dir("./adapters/http/views/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", staticFiles))
	r.Handle("/favicon.ico", staticFiles)

	t, err := views.NewTemplates()
	if err != nil {
		return nil, err
	}

	h := &Handler{game: views.NewGameView(t), currentGame: &game.Field{}}
	r.Get("/", h.show)
	r.Get("/health", h.healthCheck)
	r.Post("/new", h.newGame)
	r.Post("/reveal", h.reveal)
	r.Post("/flag", h.flag)
	r.Post("/unflag", h.unflag)
	return r, nil
}

type Handler struct {
	currentGame *game.Field
	game        *views.GameView
}

func (h *Handler) show(w http.ResponseWriter, r *http.Request) {
	h.game.Render(w, h.currentGame)
}

func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("OK"))
}

func (h *Handler) newGame(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("parse form: " + err.Error()))
		return
	}

	width, height, mines, err := extractGameParams(&r.Form)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("extract params: " + err.Error()))
		return
	}

	f, err := game.NewField(width, height, mines)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("new: " + err.Error()))
		return
	}

	h.currentGame = f
	h.game.Render(w, h.currentGame)
}

func (h *Handler) reveal(w http.ResponseWriter, r *http.Request) {
	x, y, err := parseCoordinates(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("reveal: " + err.Error()))
		return
	}

	if err := h.currentGame.Reveal(x, y); err != nil {
		h.currentGame.RevealAll()
		w.Header().Set("HX-Retarget", "grid")
		w.Header().Set("Hx-Refresh", "true")
		h.game.RenderGrid(w, h.currentGame)
		return
	}

	h.game.RenderCell(w, h.currentGame, x, y)
}

func (h *Handler) flag(w http.ResponseWriter, r *http.Request) {
	x, y, err := parseCoordinates(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("flag: " + err.Error()))
		return
	}

	h.currentGame.Flag(x, y)
	h.game.RenderCell(w, h.currentGame, x, y)
}

func (h *Handler) unflag(w http.ResponseWriter, r *http.Request) {
	x, y, err := parseCoordinates(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unflag: " + err.Error()))
		return
	}

	h.currentGame.Unflag(x, y)
	h.game.RenderCell(w, h.currentGame, x, y)
}

func parseCoordinates(r *http.Request) (int, int, error) {
	if err := r.ParseForm(); err != nil {
		return 0, 0, fmt.Errorf("parseCoordinates: parseForm: %w", err)
	}

	x, err := strconv.Atoi(r.FormValue("x"))
	if err != nil {
		return 0, 0, fmt.Errorf("parseCoordinates: x: %w", err)
	}

	y, err := strconv.Atoi(r.FormValue("y"))
	if err != nil {
		return 0, 0, fmt.Errorf("parseCoordinates: y: %w", err)
	}

	return x, y, nil
}

func extractGameParams(vv *url.Values) (int, int, int, error) {
	width, err := strconv.Atoi(vv.Get("width"))
	if err != nil {
		return 0, 0, 0, fmt.Errorf("width: %w", err)
	}

	height, err := strconv.Atoi(vv.Get("height"))
	if err != nil {
		return 0, 0, 0, fmt.Errorf("height: %w", err)
	}

	mines, err := strconv.Atoi(vv.Get("mines"))
	if err != nil {
		return 0, 0, 0, fmt.Errorf("mines: %w", err)
	}

	return width, height, mines, nil
}
