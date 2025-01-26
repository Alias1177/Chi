package middleware

import (
	"context"
	"forLessons/internal/api/auth"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Токен авторизации не найден", http.StatusUnauthorized)
			return
		}

		if len(tokenStr) > len("Bearer ") {
			tokenStr = tokenStr[len("Bearer "):]
		}

		claims, err := auth.ValidateJWT(tokenStr)
		if err != nil {
			http.Error(w, "Не правильный токен или его действие закончилось", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "username", claims.Signature)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
