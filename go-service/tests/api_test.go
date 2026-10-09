package tests

import (
	"bytes"
	"encoding/json"
	"go-service/internal/controller"
	"go-service/internal/model"
	"go-service/internal/repository"
	"go-service/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupRouter() *http.ServeMux {
	repo := repository.NewInMemoryUserRepository()
	svc := service.NewUserService(repo)
	ctrl := controller.NewUserController(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/users", ctrl.HandleUsers)
	mux.HandleFunc("/api/users/", ctrl.HandleUserByID)
	return mux
}

func TestAPIEndpoints(t *testing.T) {
	mux := setupRouter()

	// Успішний POST
	payload := []byte(`{"name": "Ілля", "email": "illia@example.com"}`)
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var created model.User
	json.NewDecoder(rr.Body).Decode(&created)
	if created.ID != 1 {
		t.Errorf("expected id 1, got %v", created.ID)
	}

	// Успішний GET
	req, _ = http.NewRequest("GET", "/api/users/1", nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// GET неіснуючого ресурсу (404)
	req, _ = http.NewRequest("GET", "/api/users/999", nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	// Успішне видалення
	req, _ = http.NewRequest("DELETE", "/api/users/1", nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code for delete: got %v want %v", status, http.StatusNoContent)
	}
}
