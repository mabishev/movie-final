package middleware

import (
	"context"
	"net/http"

	"github.com/marcokz/movie-final/internal/auth"
)

// Создаем кастомный тип для ключа контекста. Это уменьшит вероятность
// коллизий с другими значениями, которые могут быть сохранены в контексте.
type ContextKey string

const UserContextKey ContextKey = "auth_token"

func Authorize(roles ...string) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(string(UserContextKey))
			if err != nil {
				http.Error(w, "Отсутствует токен", http.StatusUnauthorized)
				return
			}

			// Проверяем JWT токен
			claims, err := auth.ValidationJWT(cookie.Value)
			if err != nil {
				http.Error(w, "Неверный токен", http.StatusUnauthorized)
				return
			}

			// Сохраняем данные о пользователе в контексте запроса
			ctx := context.WithValue(r.Context(), UserContextKey, claims)

			// Проверяем роли пользователя
			// Если хотя бы одна из ролей пользователя совпадает с ролями, переданными в параметрах, то пропускаем запрос
			// Если нет, то возвращаем ошибку
			for _, role := range roles {
				if role == claims.Role {
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}
			w.WriteHeader(http.StatusForbidden)
		}
	}
}
