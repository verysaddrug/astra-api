package service

import (
	"astra-api/internal/model"
	"astra-api/internal/repository"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo   repository.UserRepositoryInterface
	adminToken string
}

func NewAuthService(userRepo repository.UserRepositoryInterface, adminToken string) *AuthService {
	return &AuthService{userRepo: userRepo, adminToken: adminToken}
}

func (s *AuthService) Register(login, password, adminToken string) (*model.User, error) {
	if adminToken != s.adminToken {
		return nil, errors.New("invalid admin token")
	}
	if err := validateLogin(login); err != nil {
		return nil, err
	}
	if err := validatePassword(password); err != nil {
		return nil, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		ID:        uuid.New().String(),
		Login:     login,
		Password:  string(hash),
		CreatedAt: time.Now(),
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Authenticate(login, password string) (*model.User, error) {
	user, err := s.userRepo.GetByLogin(login)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}
	return user, nil
}

func validateLogin(login string) error {
	if len(login) < 8 {
		return errors.New("login must be at least 8 characters")
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+$`, login)
	if !matched {
		return errors.New("login must contain only latin letters and digits")
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	var upper, lower, digit, special bool
	for _, c := range password {
		switch {
		case 'A' <= c && c <= 'Z':
			upper = true
		case 'a' <= c && c <= 'z':
			lower = true
		case '0' <= c && c <= '9':
			digit = true
		case strings.ContainsRune("!@#$%^&*()-_=+[]{}|;:',.<>?/`~", c):
			special = true
		}
	}
	if !(upper && lower) {
		return errors.New("password must contain both upper and lower case letters")
	}
	if !digit {
		return errors.New("password must contain at least one digit")
	}
	if !special {
		return errors.New("password must contain at least one special character")
	}
	return nil
}
