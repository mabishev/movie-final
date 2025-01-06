package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

// Секретный ключ для подписи токена
var jwtKey = []byte("your_secret_key")

type Claims struct {
	ID   int64  `json:"id"`
	Role string `json:"role"`
	jwt.StandardClaims
}

func GenerateJWT(id int64, email string, role string) (string, error) {
	// Устанавливаем время жизни токена
	expirationTime := time.Now().Add(time.Hour * 24)
	claims := &Claims{
		ID:   id,
		Role: role,
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

func ValidationJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return &Claims{}, err
	}
	if !token.Valid {
		return &Claims{}, errors.New("invalid token")
	}
	return claims, nil
}
