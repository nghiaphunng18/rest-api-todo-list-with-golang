package services

import (
	"errors"
	"todo-api/internal/models"
	"todo-api/internal/repository"
	"todo-api/internal/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrUsernameExists = errors.New("username already exists")

type AuthService struct {
    userStore *repository.UserStore
    jwtSecret string
}

func NewAuthService(userStore *repository.UserStore, jwtSecret string) *AuthService {
    return &AuthService{
        userStore: userStore,
        jwtSecret: jwtSecret,
    }
}

func (s *AuthService) Register(username, password string) error {
	_, err := s.userStore.GetUserByUsername(username)
	if err == nil {
		return ErrUsernameExists
	}
    if !errors.Is(err, gorm.ErrRecordNotFound) {
        return err
    }
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    user := &models.User{
        Username: username,
        PasswordHash: string(hashed),
    }

    return s.userStore.CreateUser(user)
}

func (s *AuthService) Login(username, password string) (string, error) {
    user, err := s.userStore.GetUserByUsername(username)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return "", errors.New("invalid credentials")
        }
        return "", err
    }

    if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
        return "", errors.New("invalid credentials")
    }

    token, err := utils.GenerateJWT(user.ID, s.jwtSecret)
    if err != nil {
        return "", err
    }

    return token, nil
}
