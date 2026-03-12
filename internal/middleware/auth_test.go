package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"todo-api/internal/config"
	"todo-api/internal/middleware"
	"todo-api/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	secret := "test_secret_key"
	// setup config JWT to have a short expiration time for testing
	utils.SetJWTConfig(config.JWTConfig{Expiration: time.Hour})

	// create dummy handler to check if middleware passes the request correctly
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check user id is set in context
		userID := r.Context().Value(utils.UserIDKey)
		assert.NotNil(t, userID)
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Missing_Authorization_Header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/todos", nil)
		rr := httptest.NewRecorder()

		handler := middleware.AuthMiddleware(secret)(nextHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Invalid_Token_Format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/todos", nil)
		req.Header.Set("Authorization", "InvalidFormat token123")
		rr := httptest.NewRecorder()

		handler := middleware.AuthMiddleware(secret)(nextHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Invalid_JWT_Token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/todos", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		rr := httptest.NewRecorder()

		handler := middleware.AuthMiddleware(secret)(nextHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Valid_Token_Success", func(t *testing.T) {
		// create valid token
		userID := uint(123)
		token, _ := utils.GenerateJWT(userID, secret)

		// make request with valid token
		req := httptest.NewRequest("GET", "/todos", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		handler := middleware.AuthMiddleware(secret)(nextHandler)
		handler.ServeHTTP(rr, req)

		// check response is 200 and user id is set in context
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}