package api

import (
	"encoding/json"
	"forLessons/internal/api/auth"
	"forLessons/internal/api/jwt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) GetLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userList := make([]User, 0, len(h.Users))
	for _, user := range h.Users {
		userList = append(userList, user)
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userList); err != nil {
		http.Error(w, "Ошибка при encode ", http.StatusInternalServerError)
		return
	}
}
func (h *Handler) SetLogin(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Ошибка при decode", http.StatusBadRequest)
		return
	}
	if user.ID == "" {
		http.Error(w, "User ID обязателен", http.StatusBadRequest)
		return
	}

	if _, exists := h.Users[user.ID]; exists {
		http.Error(w, "Пользователь с таким ID уже есть", http.StatusConflict)
		return
	}

	h.Users[user.ID] = user
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Ошибка при encode", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	user, exists := h.Users[id]
	if !exists {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Ошибка при encode", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) LogPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	var updatedUser User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Не правильный формат json", http.StatusBadRequest)
		return
	}

	if _, exists := h.Users[id]; !exists {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	updatedUser.ID = id
	h.Users[id] = updatedUser
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedUser); err != nil {
		http.Error(w, "Ошибка при encode", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) LoginDel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	if _, exists := h.Users[id]; !exists {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	delete(h.Users, id)
	w.WriteHeader(http.StatusNoContent)
}
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Не правильный формат", http.StatusBadRequest)
		return
	}

	if _, exists := h.Users[user.ID]; exists {
		http.Error(w, "Пользователь уже существует", http.StatusConflict)
		return
	}

	hashedPassword, err := jwt.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Ошибка при создании пользователя", http.StatusInternalServerError)
		return
	}

	user.Password = hashedPassword
	h.Users[user.ID] = user

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Ошибка при encode", http.StatusInternalServerError)
		return
	}
}
func (h *Handler) LoginRegister(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Не правильный формат", http.StatusBadRequest)
		return
	}

	var userFound *User
	for _, user := range h.Users {
		if user.Username == credentials.Username {
			userFound = &user
			break
		}
	}

	if userFound == nil {
		http.Error(w, "Неправильный логин или пароль", http.StatusUnauthorized)
		return
	}

	if !jwt.CheckPasswordHash(credentials.Password, userFound.Password) {
		http.Error(w, "Неправильный логин или пароль", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateJWT(userFound.Username)
	if err != nil {
		http.Error(w, "Ошибка при генерации токена", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
