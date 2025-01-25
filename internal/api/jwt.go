package api

import (
	"context"
	"encoding/json"
	"forLessons/internal/api/auth"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid format", http.StatusBadRequest)
		return
	}

	if _, exists := h.Users[user.ID]; exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	user.Password = hashedPassword
	h.Users[user.ID] = user

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
func (h *Handler) LoginRegist(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid format", http.StatusBadRequest)
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
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if !CheckPasswordHash(credentials.Password, userFound.Password) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateJWT(userFound.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Authorization token missing", http.StatusUnauthorized)
			return
		}

		if len(tokenStr) > len("Bearer ") {
			tokenStr = tokenStr[len("Bearer "):]
		}

		claims, err := auth.ValidateJWT(tokenStr)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "username", claims.Signature)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
