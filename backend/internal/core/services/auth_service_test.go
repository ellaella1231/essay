package services

import (
	"context"
	"testing"

	"ai-essay-master-backend/internal/repository/memory"
)

func TestVerifyOTPCreatesUserAndToken(t *testing.T) {
	repo := memory.NewUserRepository()
	service := NewAuthService(repo)

	token, user, err := service.VerifyOTP(context.Background(), "+8613800138000", "1234")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
	if user == nil {
		t.Fatal("expected user")
	}
	if user.Phone != "+8613800138000" {
		t.Fatalf("unexpected phone: %s", user.Phone)
	}
	if user.FreeTrialCount != 3 {
		t.Fatalf("unexpected free trial count: %d", user.FreeTrialCount)
	}
}

func TestVerifyOTPRejectsBadOTP(t *testing.T) {
	repo := memory.NewUserRepository()
	service := NewAuthService(repo)

	_, _, err := service.VerifyOTP(context.Background(), "+8613800138000", "9999")
	if err == nil {
		t.Fatal("expected invalid OTP error")
	}
}

func TestVerifyOTPRejectsBadPhone(t *testing.T) {
	repo := memory.NewUserRepository()
	service := NewAuthService(repo)

	_, _, err := service.VerifyOTP(context.Background(), "13800138000", "1234")
	if err == nil {
		t.Fatal("expected invalid phone error")
	}
}
