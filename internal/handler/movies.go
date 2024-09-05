package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/marcokz/movie-final/internal/entity"
)

type MoviesRepo interface {
	CreateMovie(ctx context.Context, e entity.Movie) error
	GetMovies(ctx context.Context) ([]entity.Movie, error)
	GetMoviesByID(ctx context.Context, id int64) (entity.Movie, error)
	UpdateMovie(ctx context.Context, e entity.Movie) error
	DeleteMovieByID(ctx context.Context, id int64) error
}

type MovieHandler struct {
	moviesRepo MoviesRepo
}

func NewMovieHandler(m MoviesRepo) *MovieHandler {
	return &MovieHandler{moviesRepo: m}
}

func (h *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	var e entity.Movie

	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.moviesRepo.CreateMovie(r.Context(), e)
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

	if err := json.NewEncoder(w).Encode(movies); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *MovieHandler) GetMoviesByID(w http.ResponseWriter, r *http.Request) {
	pathValue := r.PathValue("id")
	id, err := strconv.ParseInt(pathValue, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m, err := h.moviesRepo.GetMoviesByID(r.Context(), id)
	if errors.Is(err, entity.ErrNotFound) { //???
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

type UpdateMovieRequest struct {
	Name string `json:"name"`
	Year int    `json:"year"`
}

func (h *MovieHandler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
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

	m := entity.Movie{
		ID:   id,
		Name: update.Name,
		Year: update.Year,
	}

	err = h.moviesRepo.UpdateMovie(r.Context(), m)
	if errors.Is(err, entity.ErrNotFound) { // ???? разве они совпадут?
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "movie update successfully"})
}

func (h *MovieHandler) DeleteMovieByID(w http.ResponseWriter, r *http.Request) {
	pathValue := r.PathValue("id")
	id, err := strconv.ParseInt(pathValue, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.moviesRepo.DeleteMovieByID(r.Context(), id)
	if errors.Is(err, entity.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
