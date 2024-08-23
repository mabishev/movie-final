package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/marcokz/movie-final/internal/movie"
)

type MoviesRepo interface {
	CreateMovie(ctx context.Context, movie movie.Movie) error
	GetMovies(ctx context.Context) []movie.Movie
	GetMoviesByID(ctx context.Context, id int64) (movie.Movie, error)
	UpdateMovie(ctx context.Context, m movie.Movie) error
	DeleteMovieByID(ctx context.Context, id int64) error
}

type MovieHandler struct {
	moviesRepo MoviesRepo
}

func NewMovieHandler(m MoviesRepo) *MovieHandler {
	return &MovieHandler{moviesRepo: m}
}

func (h *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json") // для показа в postman json
	var m movie.Movie

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.moviesRepo.CreateMovie(r.Context(), m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
}

func (h *MovieHandler) GetMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json") // для показа в postman json

	ms := h.moviesRepo.GetMovies(r.Context())

	if err := json.NewEncoder(w).Encode(ms); err != nil { // отправка данных
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *MovieHandler) GetMoviesByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json") // для показа в postman json
	pathValue := r.PathValue("id")
	id, err := strconv.ParseInt(pathValue, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m, err := h.moviesRepo.GetMoviesByID(r.Context(), id)
	if errors.Is(err, movie.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	if err := json.NewEncoder(w).Encode(m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

type UpdateMovieRequest struct {
	Name string `json:"name"`
	Year int    `json:"year"`
}

func (h *MovieHandler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json") // для показа в postman json
	pathValue := r.PathValue("id")
	id, err := strconv.ParseInt(pathValue, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var update UpdateMovieRequest

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m := movie.Movie{
		ID:   id,
		Name: update.Name,
		Year: update.Year,
	}

	err = h.moviesRepo.UpdateMovie(r.Context(), m)
	if errors.Is(err, movie.ErrNotFound) { // ???? разве они совпадут?
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "movie update successfully"})
}

func (h *MovieHandler) DeleteMovieByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json") // для показа в postman json
	pathValue := r.PathValue("id")
	id, err := strconv.ParseInt(pathValue, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.moviesRepo.DeleteMovieByID(r.Context(), id)
	if errors.Is(err, movie.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
