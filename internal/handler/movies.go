package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/marcokz/movie-final/internal/auth"
	"github.com/marcokz/movie-final/internal/entity"
	"github.com/marcokz/movie-final/internal/middleware"
)

type MoviesRepo interface {
	CreateMovie(ctx context.Context, m entity.Movie) error
	GetMovies(ctx context.Context) ([]entity.Movie, error)
	GetMoviesByID(ctx context.Context, id int64) (entity.Movie, error)
	UpdateMovieByID(ctx context.Context, m entity.Movie) error
	DeleteMovieByID(ctx context.Context, id int64) error
}

type MovieHandler struct {
	moviesRepo MoviesRepo
}

func NewMovieHandler(m MoviesRepo) *MovieHandler {
	return &MovieHandler{moviesRepo: m}
}

type MovieResponse struct {
	Name        string `json:"name"`
	Year        int    `json:"year"`
	Description string `json:"description"`
}

func (h *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var create MovieResponse

	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	movie := entity.Movie{
		Name:        create.Name,
		Year:        create.Year,
		Description: create.Description,
	}

	err := h.moviesRepo.CreateMovie(r.Context(), movie)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *MovieHandler) GetMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := h.moviesRepo.GetMovies(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movies)
}

func (h *MovieHandler) GetMoviesByID(w http.ResponseWriter, r *http.Request) {
	pathValue := r.PathValue("id")
	id, err := strconv.ParseInt(pathValue, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m, err := h.moviesRepo.GetMoviesByID(r.Context(), id)
	if errors.Is(err, entity.ErrMovieNotFound) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(m)
}

func (h *MovieHandler) UpdateMovieByID(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	pathValue := r.PathValue("id")
	id, err := strconv.ParseInt(pathValue, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var update MovieResponse

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m := entity.Movie{
		ID:   id,
		Name: update.Name,
		Year: update.Year,
	}

	err = h.moviesRepo.UpdateMovieByID(r.Context(), m)
	if errors.Is(err, entity.ErrMovieNotFound) { // ???? разве они совпадут?
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "movie update successfully"})
}

func (h *MovieHandler) DeleteMovieByID(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	pathValue := r.PathValue("id")
	id, err := strconv.ParseInt(pathValue, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.moviesRepo.DeleteMovieByID(r.Context(), id)
	if errors.Is(err, entity.ErrMovieNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "movie delete successfully"})
}
