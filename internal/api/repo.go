package api

import (
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	R     *chi.Mux
	Users map[string]User
}

func New() *Handler {
	return &Handler{
		R:     chi.NewRouter(),
		Users: make(map[string]User),
	}
}

type User struct {
	Password string `json:"password"`
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}
