package chain

import (
	"forLessons/internal/api"
	middleware2 "forLessons/internal/api/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func AttachHandlers(r chi.Router, handler *api.Handler) {
	if r == nil || handler == nil {
		return
	}

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Group(func(public chi.Router) {
		public.Post("/register", handler.Register)
		public.Post("/login", handler.LoginRegister)
	})
	r.Group(func(auth chi.Router) {
		auth.Use(middleware2.AuthMiddleware)
		auth.Get("/users", handler.GetLogin)
		auth.Route("/api", func(apiR chi.Router) {
			apiR.Post("/login", handler.SetLogin)
			apiR.Get("/login/{id}", handler.Login)
			apiR.Put("/login/{id}", handler.LogPut)
			apiR.Delete("/login/{id}", handler.LoginDel)
		})
	})
}
