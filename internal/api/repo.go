package api

import "github.com/go-chi/chi/v5"

type Handler struct {
	R *chi.Mux
}

func New() *Handler {
	return &Handler{
		R: chi.NewMux(),
	}
}

var users = make(map[string]User)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}
