package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func AttachHandlers(r chi.Router, handler *Handler) {
	r.Use(middleware.Recoverer)
	r.Route("/api", func(apiR chi.Router) {
		apiR.Use(middleware.Logger)
		apiR.Get("/login", handler.GetLogin)
		apiR.Post("/login", handler.SetLogin)
		apiR.Get("/login/{id}", handler.Login)
		apiR.Put("/login/{id}", handler.LogPut)
		apiR.Delete("/login/{id}", handler.LoginDel)
	})
}
