package middleware

import (
	"context"
	"net/http"

	"github.com/marcokz/movie-final/internal/auth"
)

// Создаем кастомный тип для ключа контекста. Это уменьшит вероятность
// коллизий с другими значениями, которые могут быть сохранены в контексте.
type ContextKey string

const UserContextKey ContextKey = "user"

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_cookie")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Отсутствует токен"))
			return
		}

		// Проверяем JWT токен
		claims, err := auth.ValidationJWT(cookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Неверный токен"))
			return
		}

		// Сохраняем данные о пользователе в контексте запроса
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		// Передаем управление следующему обработчику с новым контекстом
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
