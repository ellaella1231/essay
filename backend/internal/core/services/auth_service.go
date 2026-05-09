package services

import (
	"context"
	"errors"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"ai-essay-master-backend/internal/core/domain"
)

type AuthService interface {
	VerifyOTP(ctx context.Context, phone string, otp string) (string, *domain.User, error)
}

type UserRepository interface {
	FindByPhone(ctx context.Context, phone string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
}

type authService struct {
	userRepo UserRepository
}

func NewAuthService(repo UserRepository) AuthService {
	return &authService{userRepo: repo}
}

func (s *authService) VerifyOTP(ctx context.Context, phone string, otp string) (string, *domain.User, error) {
	// Validate phone +86
	match, _ := regexp.MatchString(`^\+861[3-9]\d{9}$`, phone)
	if !match {
		return "", nil, errors.New("invalid phone format, must be +86 followed by 11 digits")
	}

	// Mock OTP check
	if otp != "1234" {
		return "", nil, errors.New("invalid OTP")
	}

	user, err := s.userRepo.FindByPhone(ctx, phone)
	if err != nil {
		// If not found, create
		user = &domain.User{
			Phone:          phone,
			FreeTrialCount: 3,
		}
		err = s.userRepo.Create(ctx, user)
		if err != nil {
			return "", nil, err
		}
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "fallback_secret_for_dev_only"
	}

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}
