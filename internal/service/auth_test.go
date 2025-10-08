package service

import (
	mocksgen "astra-api/internal/mocks/gomock"
	"astra-api/internal/model"
	"errors"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)
	authService := NewAuthService(userRepo, "admin123")

	userRepo.EXPECT().Create(gomock.Any()).Return(nil)

	user, err := authService.Register("testuser123", "Password123!", "admin123")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Login != "testuser123" {
		t.Fatalf("expected login 'testuser123', got %s", user.Login)
	}
	if user.ID == "" {
		t.Fatal("expected user ID to be set")
	}
}

func TestAuthService_Register_InvalidAdminToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)
	authService := NewAuthService(userRepo, "admin123")

	_, err := authService.Register("testuser123", "Password123!", "wrongtoken")

	if err == nil {
		t.Fatal("expected error for invalid admin token")
	}
	if err.Error() != "invalid admin token" {
		t.Fatalf("expected 'invalid admin token', got %s", err.Error())
	}
}

func TestAuthService_Register_InvalidLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)
	authService := NewAuthService(userRepo, "admin123")

	testCases := []struct {
		name     string
		login    string
		expected string
	}{
		{"too short", "short", "login must be at least 8 characters"},
		{"with special chars", "user@123", "login must contain only latin letters and digits"},
		{"with spaces", "user 123", "login must contain only latin letters and digits"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := authService.Register(tc.login, "Password123!", "admin123")
			if err == nil {
				t.Fatal("expected error for invalid login")
			}
			if err.Error() != tc.expected {
				t.Fatalf("expected '%s', got '%s'", tc.expected, err.Error())
			}
		})
	}
}

func TestAuthService_Register_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)
	authService := NewAuthService(userRepo, "admin123")

	testCases := []struct {
		name     string
		password string
		expected string
	}{
		{"too short", "short", "password must be at least 8 characters"},
		{"no uppercase", "password123!", "password must contain both upper and lower case letters"},
		{"no lowercase", "PASSWORD123!", "password must contain both upper and lower case letters"},
		{"no digit", "Password!", "password must contain at least one digit"},
		{"no special", "Password123", "password must contain at least one special character"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := authService.Register("testuser123", tc.password, "admin123")
			if err == nil {
				t.Fatal("expected error for invalid password")
			}
			if err.Error() != tc.expected {
				t.Fatalf("expected '%s', got '%s'", tc.expected, err.Error())
			}
		})
	}
}

func TestAuthService_Register_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)
	authService := NewAuthService(userRepo, "admin123")

	userRepo.EXPECT().Create(gomock.Any()).Return(errors.New("database error"))

	_, err := authService.Register("testuser123", "Password123!", "admin123")

	if err == nil {
		t.Fatal("expected error from repository")
	}
	if err.Error() != "database error" {
		t.Fatalf("expected 'database error', got %s", err.Error())
	}
}

func TestAuthService_Authenticate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)
	authService := NewAuthService(userRepo, "admin123")

	// Create a real password hash for testing
	password := "Password123!"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to generate password hash: %v", err)
	}

	user := &model.User{
		ID:       uuid.New().String(),
		Login:    "testuser123",
		Password: string(hash),
	}

	userRepo.EXPECT().GetByLogin("testuser123").Return(user, nil)

	result, err := authService.Authenticate("testuser123", password)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Login != "testuser123" {
		t.Fatalf("expected login 'testuser123', got %s", result.Login)
	}
}

func TestAuthService_Authenticate_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)
	authService := NewAuthService(userRepo, "admin123")

	userRepo.EXPECT().GetByLogin("nonexistent").Return(nil, errors.New("not found"))

	_, err := authService.Authenticate("nonexistent", "Password123!")

	if err == nil {
		t.Fatal("expected error for user not found")
	}
	if err.Error() != "user not found" {
		t.Fatalf("expected 'user not found', got %s", err.Error())
	}
}

func TestAuthService_Authenticate_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocksgen.NewMockUserRepositoryInterface(ctrl)
	authService := NewAuthService(userRepo, "admin123")

	user := &model.User{
		ID:       uuid.New().String(),
		Login:    "testuser123",
		Password: "$2a$10$abcdefghijklmnopqrstuvwxyz", // bcrypt hash for "Password123!"
	}

	userRepo.EXPECT().GetByLogin("testuser123").Return(user, nil)

	_, err := authService.Authenticate("testuser123", "WrongPassword123!")

	if err == nil {
		t.Fatal("expected error for invalid password")
	}
	if err.Error() != "invalid password" {
		t.Fatalf("expected 'invalid password', got %s", err.Error())
	}
}
