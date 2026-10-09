package tests

import (
	"go-service/internal/model"
	"go-service/internal/repository"
	"go-service/internal/service"
	"testing"
)

func TestUserService(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	svc := service.NewUserService(repo)

	// Перевірка створення
	u := model.User{Name: "Ілля", Email: "illia@example.com"}
	created, err := svc.CreateUser(u)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if created.ID != 1 {
		t.Errorf("Expected ID 1, got %v", created.ID)
	}

	// Перевірка валідації (відсутній email)
	_, err = svc.CreateUser(model.User{Name: "Ілля"})
	if err != service.ErrInvalidData {
		t.Errorf("Expected ErrInvalidData, got %v", err)
	}

	// Перевірка пошуку
	found, err := svc.GetUserByID(created.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if found.Name != "Ілля" {
		t.Errorf("Expected name 'Ілля', got %v", found.Name)
	}

	// Перевірка оновлення
	updated, err := svc.UpdateUser(created.ID, model.User{Name: "Ілля Михалюк"})
	if err != nil || updated.Name != "Ілля Михалюк" {
		t.Errorf("Failed to update user")
	}

	// Перевірка видалення
	err = svc.DeleteUser(created.ID)
	if err != nil {
		t.Errorf("Expected no error on delete, got %v", err)
	}
	_, err = svc.GetUserByID(created.ID)
	if err != repository.ErrNotFound {
		t.Errorf("Expected ErrNotFound after deletion, got %v", err)
	}
}
