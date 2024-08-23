package handler

import (
	"time"

	"github.com/golang-jwt/jwt"
)

// Секретный ключ для подписи токена
var jwtKey = []byte("your_secret_key")

type Claims struct { // "claims" = "претензии, требования"
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWT(email string) (string, error) {
	// Устанавливаем время жизни токена
	expirationTime := time.Now().Add(time.Hour * 24) // "expiration time" = "срок годности/действия"
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// Создаём токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
