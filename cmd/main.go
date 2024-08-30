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

	userRepo := postgresdb.NewUserRepo(pool)

	mux := http.NewServeMux() // на каждый http запрос запускается отдельная go рутина

	withJson := middleware.WithContentTypeJSON(mux)

	adminOnly := middleware.Authorize("admin")
	userAndAdmin := middleware.Authorize("admin", "user")

	mh := handler.NewMovieHandler(movieRepo)
	mux.HandleFunc("POST /movies", adminOnly(mh.CreateMovie))
	mux.HandleFunc("GET /movies", userAndAdmin(mh.GetMovies)) //  если будет "/" в конце адреса, то обрабатываются все адреса после слэша
	mux.HandleFunc("GET /movies/{id}", userAndAdmin(mh.GetMoviesByID))
	mux.HandleFunc("PUT /movies/{id}", adminOnly(mh.DeleteMovieByID))
	mux.HandleFunc("DELETE /movies/{id}", adminOnly(mh.DeleteMovieByID))

	uh := handler.NewUserHandler(userRepo)
	mux.HandleFunc("POST /user/create", userAndAdmin(uh.CreateUser))
	mux.HandleFunc("POST /user/auth", userAndAdmin(uh.SignIn))
	mux.Handle("PUT /user/update", userAndAdmin(uh.UpdateUserInfo))

	server := &http.Server{
		Addr:    ":8080",
		Handler: withJson,
	}

	server.ListenAndServe()
}
