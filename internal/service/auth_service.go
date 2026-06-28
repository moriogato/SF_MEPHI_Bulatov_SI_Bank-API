package service

import (
	"bank-api/internal/models"
	"bank-api/internal/repository"
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(username, email, password string) (*models.User, error) {
	if !isValidEmail(email) {
		return nil, errors.New("invalid email format")
	}
	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 chars")
	}
	existing, _ := s.userRepo.FindByEmail(email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
	}
	err = s.userRepo.Create(user)
	return user, err
}

func (s *AuthService) Login(email, password string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}
