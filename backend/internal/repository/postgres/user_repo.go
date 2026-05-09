package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"ai-essay-master-backend/internal/core/domain"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(ctx, "SELECT id, phone, free_trial_count, subscription_expiry, created_at FROM users WHERE phone = $1", phone).Scan(
		&u.ID, &u.Phone, &u.FreeTrialCount, &u.SubscriptionExpiry, &u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	err := r.db.QueryRow(ctx, "INSERT INTO users (phone, free_trial_count) VALUES ($1, $2) RETURNING id, created_at", user.Phone, user.FreeTrialCount).Scan(&user.ID, &user.CreatedAt)
	return err
}
