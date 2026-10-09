package controller

import (
	"encoding/json"
	"go-service/internal/model"
	"go-service/internal/repository"
	"go-service/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type UserController struct {
	service *service.UserService
}

func NewUserController(s *service.UserService) *UserController {
	return &UserController{service: s}
}

// Універсальна функція для форматування помилок
func writeError(w http.ResponseWriter, status int, errType, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   errType,
		"message": message,
	})
}

// Обробник для /api/users (GET, POST)
func (c *UserController) HandleUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		users, _ := c.service.GetAllUsers()
		json.NewEncoder(w).Encode(users)
	case http.MethodPost:
		var u model.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			writeError(w, http.StatusBadRequest, "bad_request", "Invalid JSON")
			return
		}
		created, err := c.service.CreateUser(u)
		if err == service.ErrInvalidData {
			writeError(w, http.StatusBadRequest, "bad_request", "Missing required fields")
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(created)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
	}
}

// Обробник для /api/users/:id (GET, PUT, DELETE)
func (c *UserController) HandleUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/api/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "Invalid identifier")
		return
	}

	switch r.Method {
	case http.MethodGet:
		user, err := c.service.GetUserByID(id)
		if err == repository.ErrNotFound {
			writeError(w, http.StatusNotFound, "not_found", "User not found")
			return
		}
		json.NewEncoder(w).Encode(user)
	case http.MethodPut:
		var u model.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			writeError(w, http.StatusBadRequest, "bad_request", "Invalid JSON")
			return
		}
		updated, err := c.service.UpdateUser(id, u)
		if err == repository.ErrNotFound {
			writeError(w, http.StatusNotFound, "not_found", "User not found")
			return
		}
		json.NewEncoder(w).Encode(updated)
	case http.MethodDelete:
		err := c.service.DeleteUser(id)
		if err == repository.ErrNotFound {
			writeError(w, http.StatusNotFound, "not_found", "User not found")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
	}
}
