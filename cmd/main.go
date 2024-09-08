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

	movieRepo := postgresdb.NewMoviesRepository(pool)
	ratingRepo := postgresdb.NewRatingsRepo(pool)
	userRepo := postgresdb.NewUserRepo(pool)

	mux := http.NewServeMux() // на каждый http запрос запускается отдельная go рутина

	withJson := middleware.WithContentTypeJSON(mux)

	adminOnly := middleware.Authorize("admin")
	userAndAdmin := middleware.Authorize("admin", "user")

	m := handler.NewMovieHandler(movieRepo)
	mux.HandleFunc("POST /movies", adminOnly(m.CreateMovie))
	mux.HandleFunc("GET /movies", userAndAdmin(m.GetMovies)) //  если будет "/" в конце адреса, то обрабатываются все адреса после слэша
	mux.HandleFunc("GET /movies/{id}", userAndAdmin(m.GetMoviesByID))
	mux.HandleFunc("PUT /movies/{id}", adminOnly(m.UpdateMovieByID))
	mux.HandleFunc("DELETE /movies/{id}", adminOnly(m.DeleteMovieByID))

	u := handler.NewUserHandler(userRepo)
	mux.HandleFunc("POST /user/create", u.CreateUser)
	mux.HandleFunc("POST /user/auth", u.SignIn)
	mux.HandleFunc("GET /user/sex", userAndAdmin(u.GetUsersBySex))
	mux.HandleFunc("GET /user/age", userAndAdmin(u.GetUserByAge))
	mux.HandleFunc("PUT /user/update", userAndAdmin(u.UpdateUserInfo))

	r := handler.NewRatingsHandler(ratingRepo)
	mux.HandleFunc("GET /ratings/users/rating", r.GetUsersByRatingOfMovie)
	mux.HandleFunc("PUT /ratings/update", r.UpdateRating)

	server := &http.Server{
		Addr:    ":8080",
		Handler: withJson,
	}

	server.ListenAndServe()
}
