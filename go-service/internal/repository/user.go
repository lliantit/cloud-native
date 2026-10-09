package repository

import (
	"errors"
	"go-service/internal/model"
	"sync"
)

var ErrNotFound = errors.New("user not found")

type UserRepository interface {
	Create(user model.User) (model.User, error)
	FindByID(id int) (model.User, error)
	Update(id int, user model.User) (model.User, error)
	Delete(id int) error
	FindAll() ([]model.User, error)
}

type InMemoryUserRepository struct {
	mu     sync.RWMutex
	users  map[int]model.User
	nextID int
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:  make(map[int]model.User),
		nextID: 1,
	}
}

func (r *InMemoryUserRepository) FindAll() ([]model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []model.User
	for _, u := range r.users {
		list = append(list, u)
	}
	return list, nil
}

func (r *InMemoryUserRepository) FindByID(id int) (model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, exists := r.users[id]
	if !exists {
		return model.User{}, ErrNotFound
	}
	return user, nil
}

func (r *InMemoryUserRepository) Create(user model.User) (model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user.ID = r.nextID
	r.users[r.nextID] = user
	r.nextID++
	return user, nil
}

func (r *InMemoryUserRepository) Update(id int, user model.User) (model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, exists := r.users[id]
	if !exists {
		return model.User{}, ErrNotFound
	}
	if user.Name != "" {
		existing.Name = user.Name
	}
	if user.Email != "" {
		existing.Email = user.Email
	}
	r.users[id] = existing
	return existing, nil
}

func (r *InMemoryUserRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.users[id]; !exists {
		return ErrNotFound
	}
	delete(r.users, id)
	return nil
}
