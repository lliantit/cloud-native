package service

import (
	"errors"
	"go-service/internal/model"
	"go-service/internal/repository"
)

var ErrInvalidData = errors.New("invalid input data")

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAllUsers() ([]model.User, error) {
	return s.repo.FindAll()
}

func (s *UserService) GetUserByID(id int) (model.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) CreateUser(user model.User) (model.User, error) {
	if user.Name == "" || user.Email == "" {
		return model.User{}, ErrInvalidData
	}
	return s.repo.Create(user)
}

func (s *UserService) UpdateUser(id int, user model.User) (model.User, error) {
	return s.repo.Update(id, user)
}

func (s *UserService) DeleteUser(id int) error {
	return s.repo.Delete(id)
}
