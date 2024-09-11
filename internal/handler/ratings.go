package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/marcokz/movie-final/internal/auth"
	"github.com/marcokz/movie-final/internal/entity"
	"github.com/marcokz/movie-final/internal/middleware"
	"github.com/marcokz/movie-final/internal/postgresdb"
)

type RatingsRepo interface {
	GetUsersByRatingOfMovie(ctx context.Context, movieid, minrating, maxrating int64) ([]postgresdb.UserWithRating, error)
	UpdateRating(ctx context.Context, r entity.Rating) error
}

type RatingsHandler struct {
	ratingsRepo RatingsRepo
}

func NewRatingsHandler(r RatingsRepo) *RatingsHandler {
	return &RatingsHandler{ratingsRepo: r}
}

type GetMovieByRating struct {
	MinRating int64 `json:"minrating"`
	MaxRating int64 `json:"maxrating"`
}

func (h *RatingsHandler) GetMovieByRating(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var getMovie GetMovieByRating
	if err := json.NewDecoder(r.Body).Decode(&getMovie); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch {
	case getMovie.MinRating < 1 || getMovie.MinRating > 10 && getMovie.MaxRating < 1 || getMovie.MaxRating > 10:
		http.Error(w, "Need from 1 to 10", http.StatusBadRequest)
		return
	case getMovie.MinRating > getMovie.MaxRating:
		http.Error(w, "Min rating is higher than max", http.StatusBadRequest)
		return
	}

}

type GetUserByRatingOgMovie struct {
	MovieID   int64 `json:"movieid"`
	MinRating int64 `json:"minrating"`
	MaxRating int64 `json:"maxrating"`
}

func (h *RatingsHandler) GetUsersByRatingOfMovie(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var getUser GetUserByRatingOgMovie

	if err := json.NewDecoder(r.Body).Decode(&getUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch {
	case getUser.MinRating < 1 || getUser.MinRating > 10 && getUser.MaxRating < 1 || getUser.MaxRating > 10:
		http.Error(w, "From 1 to 10", http.StatusBadRequest)
		return
	case getUser.MinRating > getUser.MaxRating:
		http.Error(w, "Min rating is higher than max", http.StatusBadRequest)
		return
	}

	users, err := h.ratingsRepo.GetUsersByRatingOfMovie(r.Context(), getUser.MovieID, getUser.MinRating, getUser.MaxRating)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

type UpdateRating struct {
	MovieID int64 `json:"movieid"`
	Rating  int64 `json:"rating"`
}

func (h *RatingsHandler) UpdateRating(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var updateRating UpdateRating
	err := json.NewDecoder(r.Body).Decode(&updateRating)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rating := entity.Rating{
		UserId:  claims.ID,
		MovieID: updateRating.MovieID,
		Rating:  updateRating.Rating,
	}
	if rating.Rating < 1 || rating.Rating > 10 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "rating from 1 to 10"})
		return
	}

	err = h.ratingsRepo.UpdateRating(r.Context(), rating)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(map[string]string{"message": "rating update"}); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

}
