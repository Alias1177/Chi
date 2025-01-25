package main

import (
	"forLessons/internal/api"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	handlers := api.New()
	if handlers == nil {
		log.Fatal("failed to initialize handlers")
	}
	r := chi.NewRouter()
	if r == nil {
		log.Fatal("failed to initialize router")
	}
	api.AttachHandlers(r, handlers)
	log.Fatal(http.ListenAndServe(":8080", r))
}
