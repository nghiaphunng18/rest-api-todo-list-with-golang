package services_test

import (
	"testing"
	"time"

	"todo-api/internal/config"
	"todo-api/internal/domain"
	"todo-api/internal/models"
	"todo-api/internal/services"
	"todo-api/internal/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// create mock userrepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestAuthService_Login(t *testing.T) {
	// setup
	secret := "test_secret"
	utils.SetJWTConfig(config.JWTConfig{Expiration: time.Hour})
	
	t.Run("Login_Success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		authService := services.NewAuthService(mockRepo, secret)

		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		
		mockUser := &models.User{
			ID:           1,
			Username:     "testuser",
			PasswordHash: string(hashedPassword),
		}

		// define behavior of mock
		mockRepo.On("GetUserByUsername", "testuser").Return(mockUser, nil)

		// test login
		token, err := authService.Login("testuser", password)

		// check
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Login_InvalidPassword", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		authService := services.NewAuthService(mockRepo, secret)

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct_pass"), bcrypt.DefaultCost)
		mockUser := &models.User{
			Username:     "testuser",
			PasswordHash: string(hashedPassword),
		}

		mockRepo.On("GetUserByUsername", "testuser").Return(mockUser, nil)

		// test login with wrong password
		token, err := authService.Login("testuser", "wrong_pass")

		assert.Error(t, err)
		assert.Equal(t, domain.ErrBadCredentials, err)
		assert.Empty(t, token)
	})

	t.Run("Login_UserNotFound", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		authService := services.NewAuthService(mockRepo, secret)

		mockRepo.On("GetUserByUsername", "nonexistent").Return(nil, gorm.ErrRecordNotFound)

		token, err := authService.Login("nonexistent", "any_pass")

		assert.Error(t, err)
		assert.Equal(t, domain.ErrBadCredentials, err)
		assert.Empty(t, token)
	})
}

func TestAuthService_Register(t *testing.T) {
	secret := "test_secret"

	t.Run("Register_Success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		authService := services.NewAuthService(mockRepo, secret)

		mockRepo.On("GetUserByUsername", "newuser").Return(nil, gorm.ErrRecordNotFound)
		mockRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil)

		err := authService.Register("newuser", "password123")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Register_UserAlreadyExists", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		authService := services.NewAuthService(mockRepo, secret)

		mockRepo.On("GetUserByUsername", "existinguser").Return(&models.User{}, nil)

		err := authService.Register("existinguser", "password123")

		assert.Error(t, err)
		assert.Equal(t, domain.ErrUserExists, err)
	})
}