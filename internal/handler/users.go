package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/mail"
	"time"

	"github.com/marcokz/movie-final/internal/auth"
	"github.com/marcokz/movie-final/internal/entity"
	"github.com/marcokz/movie-final/internal/middleware"

	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	CreateUser(ctx context.Context, u entity.User) error
	GetUserByEmail(ctx context.Context, loginOrEmail string) (entity.User, error)
	GetUserByAge(ctx context.Context, minAge, maxAge int64) ([]entity.User, error)
	GetUserByCountry(ctx context.Context, country string) ([]entity.User, error)
	GetUserByCity(ctx context.Context, city string) ([]entity.User, error)
	GetUserBySex(ctx context.Context, sex string) ([]entity.User, error)
	UpdateUserInfo(ctx context.Context, u entity.User) error
}

type UserHandler struct {
	userRepo UserRepo
}

func NewUserHandler(u UserRepo) *UserHandler {
	return &UserHandler{userRepo: u}
}

type CreateUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID          int64
	Name        string
	Surname     string
	Sex         string
	DateOfBirth string
	Country     string
	City        string
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var create CreateUser

	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	e, err := mail.ParseAddress(create.Email)
	if err != nil {
		http.Error(w, "Incorrect email", http.StatusBadRequest)
		return
	}

	u := entity.User{
		Email:    e.Address,
		Password: create.Password,
	}

	err = h.userRepo.CreateUser(context.Background(), u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "user create successfully"})
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var request SignInRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := h.userRepo.GetUserByEmail(context.Background(), request.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error:": err.Error()})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(request.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error:": "invalid password"})
		return
	}

	tokenString, err := auth.GenerateJWT(u.ID, u.Email, "user")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to generate token"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24),
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Включайте только если используете HTTPS
	})

	w.WriteHeader(http.StatusOK)
}

type Age struct {
	MinAge int64 `json:"minage"`
	MaxAge int64 `json:"maxage"`
}

func (h *UserHandler) GetUserByAge(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var age Age

	if err := json.NewDecoder(r.Body).Decode(&age); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := h.userRepo.GetUserByAge(r.Context(), age.MinAge, age.MaxAge)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userResp := make([]User, 0, len(users))

	for _, user := range users {
		u := User{
			ID:          user.ID,
			Name:        user.Name,
			Surname:     user.Surname,
			Sex:         user.Sex,
			DateOfBirth: user.DateOfBirth.Format("2006-01-02"),
			Country:     user.Country,
			City:        user.City,
		}
		userResp = append(userResp, u)

	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userResp)
}

type GetByCountry struct {
	Country string `json:"country"`
}

func (h *UserHandler) GetUserByCountry(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var country GetByCountry
	if err := json.NewDecoder(r.Body).Decode(&country); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := h.userRepo.GetUserByCountry(r.Context(), country.Country)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userResp := make([]User, 0, len(users))

	for _, user := range users {
		u := User{
			ID:          user.ID,
			Name:        user.Name,
			Surname:     user.Surname,
			Sex:         user.Sex,
			DateOfBirth: user.DateOfBirth.Format("2006-01-02"),
			Country:     user.Country,
			City:        user.City,
		}
		userResp = append(userResp, u)

	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userResp)
}

type GetCity struct {
	City string `json:"city"`
}

func (h *UserHandler) GetUserByCity(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var c GetCity
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := h.userRepo.GetUserByCity(r.Context(), c.City)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userResp := make([]User, 0, len(users))

	for _, u := range users {
		city := User{
			ID:          u.ID,
			Name:        u.Name,
			Surname:     u.Surname,
			Sex:         u.Sex,
			DateOfBirth: u.DateOfBirth.Format("2006-01-02"),
			Country:     u.Country,
			City:        u.City,
		}
		userResp = append(userResp, city)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userResp)
}

type UserBySex struct {
	Sex string `json:"sex"`
}

func (h *UserHandler) GetUserBySex(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var userBySex UserBySex

	if err := json.NewDecoder(r.Body).Decode(&userBySex); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := h.userRepo.GetUserBySex(r.Context(), userBySex.Sex)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userResp := make([]User, 0, len(users))

	for _, u := range users {
		city := User{
			ID:          u.ID,
			Name:        u.Name,
			Surname:     u.Surname,
			Sex:         u.Sex,
			DateOfBirth: u.DateOfBirth.Format("2006-01-02"),
			Country:     u.Country,
			City:        u.City,
		}
		userResp = append(userResp, city)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userResp)
}

func (h *UserHandler) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var update User

	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	parsedDate, err := time.Parse("2006-01-02", update.DateOfBirth)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u := entity.User{
		ID:          claims.ID,
		Name:        update.Name,
		Surname:     update.Surname,
		Sex:         update.Sex,
		DateOfBirth: parsedDate,
		Country:     update.Country,
		City:        update.City,
	}

	err = h.userRepo.UpdateUserInfo(r.Context(), u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "user update successfully"})
}
