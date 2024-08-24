package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/marcokz/movie-final/internal/users"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	CreateUser(ctx context.Context, u users.User) error
	GetUserByEmail(ctx context.Context, loginOrEmail string) (users.User, error)
	UpdateUserInfo(ctx context.Context, u users.User) error
}

type UserHandler struct {
	userRepo UserRepo
}

func NewUserHandler(u UserRepo) *UserHandler {
	return &UserHandler{userRepo: u}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var u users.User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.userRepo.CreateUser(context.Background(), u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
}

type SignInRequest struct {
	Email    string `json:email`
	Password string `json:password`
}

func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

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
	tokenString, err := GenerateJWT(u.ID, u.Email) // что дальше делать с токеном?
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to generate token"})
		return
	}

	// Возвращаем JWT токен клиенту
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString}) //???
}

type UpdateUserInfo struct {
	Sex         string `json:sex`
	DateOfBirth string `json:dateofbirth`
	Country     string `json:country`
	City        string `json:country`
}

func (h *UserHandler) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	// Получение токена из заголовка Authorization
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing token"})
		return
	}

	// Удаление префикса "Bearer " из строки токена
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	// Проверка токена
	claims, err := ValidationJWT(tokenStr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid token"})
		return
	}

	var update UpdateUserInfo

	err = json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u := users.User{
		ID:          claims.ID,
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
