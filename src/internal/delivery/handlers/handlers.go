package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"users-api/src/internal/domain"
	"users-api/src/internal/errors"
)

type UserService interface {
	CreateUser(user *domain.User) error
	GetUser(id int64) (*domain.User, error)
	UpdateUser(user *domain.User) error
	DeleteUser(id int64) error
}

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeError(w, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	if err := h.userService.CreateUser(&user); err != nil {
		switch err {
		case errors.ErrInvalidInput:
			writeError(w, http.StatusBadRequest, err, "Invalid input data")
		case errors.ErrInvalidEmail:
			writeError(w, http.StatusBadRequest, err, "Invalid email format")
		default:
			writeError(w, http.StatusInternalServerError, err, "Failed to create user")
		}
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err, "Invalid user ID")
		return
	}

	user, err := h.userService.GetUser(id)
	if err != nil {
		if err == errors.ErrUserNotFound {
			writeError(w, http.StatusNotFound, err, "User not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err, "Failed to get user")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeError(w, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	if err := h.userService.UpdateUser(&user); err != nil {
		switch err {
		case errors.ErrInvalidInput:
			writeError(w, http.StatusBadRequest, err, "Invalid input data")
		case errors.ErrUserNotFound:
			writeError(w, http.StatusNotFound, err, "User not found")
		case errors.ErrInvalidEmail:
			writeError(w, http.StatusBadRequest, err, "Invalid email format")
		default:
			writeError(w, http.StatusInternalServerError, err, "Failed to update user")
		}
		return
	}

	writeJSON(w, http.StatusOK, user)
}
