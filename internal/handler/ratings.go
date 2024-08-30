package handler

import (
	"context"
	"encoding/json"
	"net/http"
)

type RatingsRepo interface {
	AddRating(ctx context.Context, userid, movieid, rating string) error
}

type RatingsHandler struct {
	RatingsRepo
}

func NewRatingsHandler(r RatingsRepo) *RatingsHandler {
	return &RatingsHandler{RatingsRepo: r}
}

type Rating int

func (h *RatingsHandler) AddRating(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var rating Rating

	err := json.NewDecoder(r.Body).Decode(&rating)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.RatingsRepo.AddRating(ctx.)
}
