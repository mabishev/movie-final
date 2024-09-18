package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/marcokz/movie-final/internal/auth"
	"github.com/marcokz/movie-final/internal/entity"
	"github.com/marcokz/movie-final/internal/middleware"
)

type RatingsRepo interface {
	GetUsersByRatingOfMovie(ctx context.Context, movieid, minrating, maxrating int64) ([]entity.UserWithRating, error)
	UpdateRating(ctx context.Context, r entity.Rating) error
}

type RatingsHandler struct {
	ratingsRepo RatingsRepo
}

func NewRatingsHandler(r RatingsRepo) *RatingsHandler {
	return &RatingsHandler{ratingsRepo: r}
}

type GetMovieByRating struct {
	UserID    int64 `json:"userid"`
	MinRating int64 `json:"minrating"`
	MaxRating int64 `json:"maxrating"`
}

func CorrectMinMaxRating(minrating, maxrating int64) (string, int) {
	switch {
	case minrating < 1 || minrating > 10:
		return "MinRating from 1 to 10", http.StatusBadRequest
	case maxrating < 1 || maxrating > 10:
		return "MaxRating from 1 to 10", http.StatusBadRequest
	case minrating > maxrating:
		return "Min rating is higher than max", http.StatusBadRequest
	}
	return "", 0
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

	str, codeInt := CorrectMinMaxRating(getMovie.MinRating, getMovie.MaxRating)
	if str != "" && codeInt != 0 {
		http.Error(w, str, codeInt)
	}

	h.ratingsRepo.GetMovieByRating
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

	str, codeInt := CorrectMinMaxRating(getUser.MinRating, getUser.MaxRating)
	if str != "" && codeInt != 0 {
		http.Error(w, str, codeInt)
	}

	users, err := h.ratingsRepo.GetUsersByRatingOfMovie(r.Context(), getUser.MovieID, getUser.MinRating, getUser.MaxRating)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
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
	json.NewEncoder(w).Encode(map[string]string{"message": "rating update"})
}
