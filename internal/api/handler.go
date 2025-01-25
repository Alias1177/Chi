package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handler) GetLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userList := make([]User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userList)
}

func (h *Handler) SetLogin(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}
	if user.ID == "" {
		http.Error(w, "ID пользователя обязателен", http.StatusBadRequest)
		return
	}

	if _, exists := users[user.ID]; exists {
		http.Error(w, "Пользователь с таким ID уже существует", http.StatusConflict)
		return
	}

	users[user.ID] = user
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	user, exists := users[id]
	if !exists {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) LogPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	var updatedUser User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	if _, exists := users[id]; !exists {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	updatedUser.ID = id
	users[id] = updatedUser
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedUser)
}
func (h *Handler) LoginDel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	if _, exists := users[id]; !exists {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	delete(users, id)
	w.WriteHeader(http.StatusNoContent)
}
