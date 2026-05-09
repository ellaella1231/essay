package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"ai-essay-master-backend/internal/core/domain"
)

type UserRepo struct {
	mu           sync.Mutex
	nextID       int
	usersByPhone map[string]*domain.User
}

func NewUserRepository() *UserRepo {
	return &UserRepo{
		nextID:       1,
		usersByPhone: make(map[string]*domain.User),
	}
}

func (r *UserRepo) FindByPhone(_ context.Context, phone string) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.usersByPhone[phone]
	if !ok {
		return nil, errors.New("user not found")
	}

	copy := *user
	return &copy, nil
}

func (r *UserRepo) Create(_ context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.usersByPhone[user.Phone]; exists {
		return errors.New("user already exists")
	}

	user.ID = r.nextID
	user.CreatedAt = time.Now()
	r.nextID++

	copy := *user
	r.usersByPhone[user.Phone] = &copy
	return nil
}
