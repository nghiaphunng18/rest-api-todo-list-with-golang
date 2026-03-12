package services_test

import (
	"testing"

	"todo-api/internal/domain"
	"todo-api/internal/models"
	"todo-api/internal/repository"
	"todo-api/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// create mock todorepository
type MockTodoRepository struct {
	mock.Mock
}

func (m *MockTodoRepository) GetTodosByUserID(userID uint, limit int, offset int) (*repository.PaginatedTodos, error) {
	args := m.Called(userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedTodos), args.Error(1)
}

func (m *MockTodoRepository) CreateTodo(todo *models.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *MockTodoRepository) GetTodoByID(id uint, userID uint) (*models.Todo, error) {
	args := m.Called(id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Todo), args.Error(1)
}

func (m *MockTodoRepository) UpdateTodo(todo *models.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *MockTodoRepository) DeleteTodo(id uint, userID uint) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

// Case Test
func TestTodoService_GetTodos(t *testing.T) {
	t.Run("GetTodos_Pagination_Logic", func(t *testing.T) {
		mockRepo := new(MockTodoRepository)
		service := services.NewTodoService(mockRepo)

		userID := uint(1)
		limit := 10
		page := 2
		expectedOffset := 10 // (2-1) * 10

		expectedData := &repository.PaginatedTodos{
			Todos:      []models.Todo{{ID: 1, Title: "Test Todo"}},
			TotalCount: 1,
		}

		// check service calculate offset correctly and call repo with correct params
		mockRepo.On("GetTodosByUserID", userID, limit, expectedOffset).Return(expectedData, nil)

		result, err := service.GetTodos(userID, limit, page)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.TotalCount)
		mockRepo.AssertExpectations(t)
	})
}

func TestTodoService_UpdateTodo(t *testing.T) {
	t.Run("UpdateTodo_Partial_Fields", func(t *testing.T) {
		mockRepo := new(MockTodoRepository)
		service := services.NewTodoService(mockRepo)

		userID := uint(1)
		todoID := uint(10)
		
		existingTodo := &models.Todo{
			ID: todoID, UserID: userID, Title: "Old Title", Completed: false,
		}

		// get current state of todo
		mockRepo.On("GetTodoByID", todoID, userID).Return(existingTodo, nil)

		// only update title, leave completed unchanged
		newTitle := "New Title"
		mockRepo.On("UpdateTodo", mock.MatchedBy(func(todo *models.Todo) bool {
			return todo.Title == newTitle && todo.Completed == false
		})).Return(nil)

		updatedTodo, err := service.UpdateTodo(userID, todoID, &newTitle, nil, nil)

		assert.NoError(t, err)
		assert.Equal(t, newTitle, updatedTodo.Title)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateTodo_NotFound", func(t *testing.T) {
		mockRepo := new(MockTodoRepository)
		service := services.NewTodoService(mockRepo)

		mockRepo.On("GetTodoByID", uint(99), uint(1)).Return(nil, gorm.ErrRecordNotFound)

		res, err := service.UpdateTodo(1, 99, nil, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, domain.ErrTodoNotFound, err)
	})
}

func TestTodoService_DeleteTodo(t *testing.T) {
	t.Run("DeleteTodo_Success", func(t *testing.T) {
		mockRepo := new(MockTodoRepository)
		service := services.NewTodoService(mockRepo)

		todoID := uint(5)
		userID := uint(1)

		mockRepo.On("GetTodoByID", todoID, userID).Return(&models.Todo{ID: todoID}, nil)
		mockRepo.On("DeleteTodo", todoID, userID).Return(nil)

		err := service.DeleteTodo(userID, todoID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}