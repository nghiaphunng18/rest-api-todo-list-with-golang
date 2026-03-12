package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"todo-api/internal/domain"
	"todo-api/internal/handlers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mock authservice
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(username, password string) error {
	args := m.Called(username, password)
	return args.Error(0)
}

func (m *MockAuthService) Login(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

// test case
func TestAuthHandler(t *testing.T) {
	t.Run("Register_Success", func(t *testing.T) {
		mockService := new(MockAuthService)
		handler := handlers.NewAuthHandler(mockService)

		input := handlers.AuthRequest{
			Username: "newuser",
			Password: "password123",
		}
		body, _ := json.Marshal(input)

		req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		mockService.On("Register", input.Username, input.Password).Return(nil).Once()

		handler.Register(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Register_ValidationError_ShortPassword", func(t *testing.T) {
		mockService := new(MockAuthService)
		handler := handlers.NewAuthHandler(mockService)

		input := handlers.AuthRequest{
			Username: "newuser",
			Password: "123", // Not: min=6
		}
		body, _ := json.Marshal(input)

		req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.Register(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockService.AssertNotCalled(t, "Register", mock.Anything, mock.Anything)
	})

	t.Run("Login_Success", func(t *testing.T) {
		mockService := new(MockAuthService)
		handler := handlers.NewAuthHandler(mockService)

		input := handlers.AuthRequest{
			Username: "testuser",
			Password: "password123",
		}
		body, _ := json.Marshal(input)
		expectedToken := "fake-jwt-token"

		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		mockService.On("Login", input.Username, input.Password).Return(expectedToken, nil).Once()

		handler.Login(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var resp map[string]string
		json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.Equal(t, expectedToken, resp["token"])
	})

	t.Run("Login_Failure_BadCredentials", func(t *testing.T) {
		mockService := new(MockAuthService)
		handler := handlers.NewAuthHandler(mockService)

		input := handlers.AuthRequest{
			Username: "testuser",
			Password: "wrongpassword",
		}
		body, _ := json.Marshal(input)

		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		mockService.On("Login", input.Username, input.Password).
			Return("", domain.ErrBadCredentials).Once()

		handler.Login(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}