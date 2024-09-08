package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/marcokz/movie-final/internal/auth"
	"github.com/marcokz/movie-final/internal/entity"
	"github.com/marcokz/movie-final/internal/middleware"

	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	CreateUser(ctx context.Context, u entity.User) error
	GetUsersByAge(ctx context.Context, minAge, maxAge int64) ([]entity.User, error)
	GetUserByCountry(ctx context.Context, country string) ([]entity.User, error)
	GetUserByEmail(ctx context.Context, loginOrEmail string) (entity.User, error)
	GetUsersBySex(ctx context.Context, sex string) ([]entity.User, error)
	UpdateUserInfo(ctx context.Context, u entity.User) error
}

type UserHandler struct {
	userRepo UserRepo
}

func NewUserHandler(u UserRepo) *UserHandler {
	return &UserHandler{userRepo: u}
}

type CreateUserUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var create CreateUserUser

	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u := entity.User{
		Email:    create.Email,
		Password: create.Password,
	}

	err := h.userRepo.CreateUser(context.Background(), u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
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

	// Сравнение пароля с хэшированным паролем
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(request.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error:": "invalid password"})
		return
	}

	//генерация jwt токена
	tokenString, err := auth.GenerateJWT(u.ID, u.Email, "user")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to generate token"})
		return
	}

	// Установка куки с JWT токеном
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

	u, err := h.userRepo.GetUsersByAge(r.Context(), age.MinAge, age.MaxAge)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
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

	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(users); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type UserBySex struct {
	Sex string `json:"sex"`
}

func (h *UserHandler) GetUsersBySex(w http.ResponseWriter, r *http.Request) {
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

	users, err := h.userRepo.GetUsersBySex(r.Context(), userBySex.Sex)
	if err != nil {
		//http.Error(w, "Failed to get users", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(users); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

type UpdateUserInfo struct {
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	Sex         string    `json:"sex"`
	DateOfBirth time.Time `json:"dateofbirth"`
	Country     string    `json:"country"`
	City        string    `json:"city"`
}

func (u *UpdateUserInfo) UnmarshalJSON(data []byte) error {
	type Alias UpdateUserInfo
	aux := &struct {
		DateOfBirth string `json:"dateofbirth"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse the date from the string format
	parsedDate, err := time.Parse("2006-01-02", aux.DateOfBirth)
	if err != nil {
		return err
	}

	u.DateOfBirth = parsedDate
	return nil
}

func (u *UpdateUserInfo) MarshalJSON() ([]byte, error) { // не работает?
	type Alias UpdateUserInfo
	return json.Marshal(&struct {
		DateOfBirth string `json:"dateofbirth"`
		*Alias
	}{
		DateOfBirth: u.DateOfBirth.Format("2006-01-02"),
		Alias:       (*Alias)(u),
	})
}

func (h *UserHandler) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	// Извлечение данных о пользователе из контекста
	claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var update UpdateUserInfo

	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		log.Println("00")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u := entity.User{
		ID:          claims.ID,
		Name:        update.Name,
		Surname:     update.Surname,
		Sex:         update.Sex,
		DateOfBirth: update.DateOfBirth,
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
