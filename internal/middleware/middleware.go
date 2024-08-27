package middleware

import (
	"context"
	"net/http"

	"github.com/marcokz/movie-final/internal/handler"
)

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_cookie")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Отсутствует токен"))
			return
		}

		// Проверяем JWT токен
		claims, err := handler.ValidationJWT(cookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Неверный токен"))
			return
		}

		// Сохраняем данные о пользователе в контексте запроса
		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
