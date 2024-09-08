package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/marcokz/movie-final/internal/entity"
)

type RatingsRepo interface {
	GetUsersByRating(ctx context.Context, movieid, minrating, maxrating int64) ([]entity.User, error)
	UpdateRating(ctx context.Context, r entity.Rating) error
}

type RatingsHandler struct {
	ratingsRepo RatingsRepo
}

func NewRatingsHandler(r RatingsRepo) *RatingsHandler {
	return &RatingsHandler{ratingsRepo: r}
}

type GetUser struct {
	MovieID   int64 `json:"movieid"`
	MinRating int64 `json:"minrating"`
	MaxRating int64 `json:"maxrating"`
}

func (h *RatingsHandler) GetUsersByRatingOfMovie(w http.ResponseWriter, r *http.Request) {
	var getUser GetUser

	err := json.NewDecoder(r.Body).Decode(&getUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := h.ratingsRepo.GetUsersByRating(r.Context(), getUser.MovieID, getUser.MinRating, getUser.MaxRating)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(users); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

type UpdateRating struct {
	UserID  int64 `json:"userid"`
	MovieID int64 `json:"movieid"`
	Rating  int64 `json:"rating"`
}

func (h *RatingsHandler) UpdateRating(w http.ResponseWriter, r *http.Request) {
	var updateRating UpdateRating
	err := json.NewDecoder(r.Body).Decode(&updateRating)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rating := entity.Rating{
		UserId:  updateRating.UserID,
		MovieID: updateRating.MovieID,
		Rating:  updateRating.Rating,
	}
	if rating.Rating < 1 || rating.Rating > 10 { //?? durys pa? kerek pe?
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "rating from 1 to 10"})
		return
	}

	err = h.ratingsRepo.UpdateRating(r.Context(), rating) //?? context
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "rating update"})
}
