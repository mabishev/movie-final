package middleware

import "net/http"

func WithContentTypeJSON(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json") //for all handlers
		next.ServeHTTP(w, r)
	})
}
