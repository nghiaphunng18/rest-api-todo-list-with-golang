package services

import (
	"errors"
	"todo-api/internal/domain"
	"todo-api/internal/models"
	"todo-api/internal/repository"
	"todo-api/internal/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IAuthService interface {
	Register(username, password string) error
	Login(username, password string) (string, error)
}

type AuthService struct {
	userRepo  repository.IUserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.IUserRepository, jwtSecret string) IAuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(username, password string) error {
	_, err := s.userRepo.GetUserByUsername(username)
	if err == nil {
		return domain.ErrUserExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ErrInternalServer
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.ErrInternalServer
	}

	user := &models.User{
		Username:     username,
		PasswordHash: string(hashed),
	}

	return s.userRepo.CreateUser(user)
}

func (s *AuthService) Login(username, password string) (string, error) {
    user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domain.ErrBadCredentials
		}
		return "", domain.ErrInternalServer
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", domain.ErrBadCredentials
	}

	return utils.GenerateJWT(user.ID, s.jwtSecret)
}
