package repository

import (
	"go-backend/models"
	"sync"
)

type UserRepository struct {
	users map[int]models.User
	mu    sync.RWMutex
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[int]models.User),
	}
}

func (r *UserRepository) GetUserByID(id int) (models.User, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user, ok := r.users[id]
	return user, ok
}

func (r *UserRepository) Save(user models.User) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
}
