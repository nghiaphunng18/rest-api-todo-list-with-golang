package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"todo-api/internal/handlers"
	"todo-api/internal/models"
	"todo-api/internal/repository"
	"todo-api/internal/utils"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mock todoservice
type MockTodoService struct {
	mock.Mock
}

func (m *MockTodoService) GetTodos(userID uint, limit int, page int) (*repository.PaginatedTodos, error) {
	args := m.Called(userID, limit, page)
	return args.Get(0).(*repository.PaginatedTodos), args.Error(1)
}

func (m *MockTodoService) CreateTodo(userID uint, title string, description string) (*models.Todo, error) {
	args := m.Called(userID, title, description)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Todo), args.Error(1)
}

func (m *MockTodoService) GetTodoByID(id uint, userID uint) (*models.Todo, error) {
	args := m.Called(id, userID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Todo), args.Error(1)
}

func (m *MockTodoService) UpdateTodo(userID uint, todoID uint, title *string, description *string, completed *bool) (*models.Todo, error) {
	args := m.Called(userID, todoID, title, description, completed)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Todo), args.Error(1)
}

func (m *MockTodoService) DeleteTodo(userID uint, todoID uint) error {
	args := m.Called(userID, todoID)
	return args.Error(0)
}

// case test
func TestTodoHandler(t *testing.T) {
	mockService := new(MockTodoService)
	handler := handlers.NewTodoHandler(mockService)
	userID := uint(1)

	t.Run("CreateTodo_Success", func(t *testing.T) {
		input := handlers.CreateTodoRequest{
			Title:       "Học Testing",
			Description: "Viết Unit Test cho Go",
		}
		body, _ := json.Marshal(input)

		req := httptest.NewRequest("POST", "/todos", bytes.NewBuffer(body))
		// user id in context
		ctx := context.WithValue(req.Context(), utils.UserIDKey, userID)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		mockService.On("CreateTodo", userID, input.Title, input.Description).
			Return(&models.Todo{ID: 1, Title: input.Title}, nil).Once()

		handler.CreateTodo(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("CreateTodo_ValidationError", func(t *testing.T) {
		// title short (< 3 char), trigger validator
		input := handlers.CreateTodoRequest{Title: "Go"} 
		body, _ := json.Marshal(input)

		req := httptest.NewRequest("POST", "/todos", bytes.NewBuffer(body))
		ctx := context.WithValue(req.Context(), utils.UserIDKey, userID)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.CreateTodo(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("GetTodoByID_Success", func(t *testing.T) {
		todoID := uint(10)
		req := httptest.NewRequest("GET", "/todos/10", nil)
		
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "10")
		
		ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
		ctx = context.WithValue(ctx, utils.UserIDKey, userID)
		req = req.WithContext(ctx)
		
		rr := httptest.NewRecorder()

		mockService.On("GetTodoByID", todoID, userID).
			Return(&models.Todo{ID: todoID, Title: "Test Todo"}, nil).Once()

		handler.GetTodoByID(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var response models.Todo
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, todoID, response.ID)
	})

	t.Run("UpdateTodo_Success", func(t *testing.T) {
		todoID := uint(10)
		newTitle := "Updated Title"
		input := handlers.UpdateTodoRequest{Title: &newTitle}
		body, _ := json.Marshal(input)

		req := httptest.NewRequest("PUT", "/todos/10", bytes.NewBuffer(body))
		
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "10")
		ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
		ctx = context.WithValue(ctx, utils.UserIDKey, userID)
		req = req.WithContext(ctx)
		
		rr := httptest.NewRecorder()

		mockService.On("UpdateTodo", userID, todoID, &newTitle, mock.Anything, mock.Anything).
			Return(&models.Todo{ID: todoID, Title: newTitle}, nil).Once()

		handler.UpdateTodo(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockService.AssertExpectations(t)
	})
}