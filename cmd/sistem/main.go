package main

import (
	"forLessons/internal/api"
	"forLessons/internal/api/chain"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	handlers := api.New()
	if handlers == nil {
		log.Fatal("Ошибка,он пустой")
	}
	r := chi.NewRouter()
	if r == nil {
		log.Fatal("Роутер моросит")
	}
	chain.AttachHandlers(r, handlers)
	log.Fatal(http.ListenAndServe(":8080", r))
}
