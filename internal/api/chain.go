package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func AttachHandlers(r chi.Router, handler *Handler) {
	if r == nil || handler == nil {
		return
	}

	r.Use(middleware.Recoverer)

	r.Group(func(public chi.Router) {
		public.Post("/register", handler.Register)
		public.Post("/login", handler.LoginRegist)
	})

	r.Group(func(auth chi.Router) {
		auth.Use(AuthMiddleware)
		auth.Get("/users", handler.GetLogin)

		auth.Route("/api", func(apiR chi.Router) {
			apiR.Use(middleware.Logger)
			apiR.Get("/login", handler.GetLogin)
			apiR.Post("/login", handler.SetLogin)
			apiR.Get("/login/{id}", handler.Login)
			apiR.Put("/login/{id}", handler.LogPut)
			apiR.Delete("/login/{id}", handler.LoginDel)
		})
	})
}
