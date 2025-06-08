package main

import (
	"log"
	"net/http"

	"github.com/marcokz/movie-final/internal/handler"
	"github.com/marcokz/movie-final/internal/middleware"
	"github.com/marcokz/movie-final/internal/postgresdb"
)

const connString = "postgres://postgres:nolan@127.0.0.1:5432"

func main() {
	pool, err := postgresdb.Connect(connString)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer pool.Close()

	movieRepo := postgresdb.NewMoviesRepo(pool)
	ratingRepo := postgresdb.NewRatingsRepo(pool)
	userRepo := postgresdb.NewUserRepo(pool)

	mux := http.NewServeMux() // на каждый http запрос запускается отдельная go рутина

	withJson := middleware.WithContentTypeJSON(mux)

	adminOnly := middleware.Authorize("admin")
	userAndAdmin := middleware.Authorize("admin", "user")

	m := handler.NewMovieHandler(movieRepo)
	mux.HandleFunc("POST /movies", adminOnly(m.CreateMovie))
	mux.HandleFunc("GET /movies", m.GetMovies)
	mux.HandleFunc("GET /movies/{id}", m.GetMoviesByID)
	mux.HandleFunc("PUT /movies/{id}", adminOnly(m.UpdateMovieByID))
	mux.HandleFunc("DELETE /movies/{id}", adminOnly(m.DeleteMovieByID))

	u := handler.NewUserHandler(userRepo)
	mux.HandleFunc("POST /user/create", u.CreateUser)
	mux.HandleFunc("POST /user/auth", u.Login)
	mux.HandleFunc("POST /user/logout", u.Logout)
	mux.HandleFunc("GET /user/age", userAndAdmin(u.GetUserByAge))
	mux.HandleFunc("GET /user/country", userAndAdmin(u.GetUserByCountry))
	mux.HandleFunc("GET /user/city", userAndAdmin(u.GetUserByCity))
	mux.HandleFunc("GET /user/sex", userAndAdmin(u.GetUserBySex))
	mux.HandleFunc("PUT /user/update", userAndAdmin(u.UpdateUserInfo))

	r := handler.NewRatingsHandler(ratingRepo)
	mux.HandleFunc("GET /ratings/movies/rating", userAndAdmin(r.GetMoviesWithRatingFromUser))
	mux.HandleFunc("GET /ratings/users/rating", userAndAdmin(r.GetUsersByRatingOfMovie))
	mux.HandleFunc("PUT /ratings/update", userAndAdmin(r.UpdateRating))

	server := &http.Server{
		Addr:    ":8080",
		Handler: withJson,
	}

	server.ListenAndServe()
}
