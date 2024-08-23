package main

import (
	"log"
	"net/http"

	"github.com/marcokz/movie-final/internal/handler"
	"github.com/marcokz/movie-final/internal/postgresdb"
)

const connString = "postgres://postgres:nolan@127.0.0.1:5432"

func main() {
	pool, err := postgresdb.Connect(connString)
	if err != nil {    // github
		log.Fatal(err)
		return
	}
	defer pool.Close()

	movieRepo := postgresdb.NewMoviesRepository(pool)

	userRepo := postgresdb.NewUserRepo(pool)

	mux := http.NewServeMux() // на каждый http запрос запускается отдельная go рутина

	mh := handler.NewMovieHandler(movieRepo)
	mux.HandleFunc("POST /movies", mh.CreateMovie)
	mux.HandleFunc("GET /movies", mh.GetMovies) //  если будет "/" в конце адреса, то обрабатываются все адреса после слэша
	mux.HandleFunc("GET /movies/{id}", mh.GetMoviesByID)
	mux.HandleFunc("PUT /movies/{id}", mh.UpdateMovie)
	mux.HandleFunc("DELETE /movies/{id}", mh.DeleteMovieByID)

	uh := handler.NewUserHandler(userRepo)
	mux.HandleFunc("POST /users", uh.CreateUser)
	mux.HandleFunc("POST /auth", uh.SignIn)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
