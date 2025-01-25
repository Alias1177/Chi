package main

import (
	"forLessons/internal/api"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	handlers := api.New()
	r := chi.NewRouter()
	api.AttachHandlers(r, handlers)
	log.Fatal(http.ListenAndServe(":8080", r))

}
