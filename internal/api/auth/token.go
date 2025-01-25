package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtSecret = []byte("your_secret_key")

func GenerateJWT(username string) (string, error) {
	claims := &jwt.MapClaims{
		"authorized": true,
		"username":   username,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
}
