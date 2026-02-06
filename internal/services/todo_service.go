package services

import (
	"errors"
	"todo-api/internal/domain"
	"todo-api/internal/models"
	"todo-api/internal/repository"

	"gorm.io/gorm"
)

type ITodoService interface {
	GetTodos(userID uint, limit int, page int) (*repository.PaginatedTodos, error)
	CreateTodo(userID uint, title string, description string) (*models.Todo, error)
	GetTodoByID(id uint, userID uint) (*models.Todo, error)
	UpdateTodo(userID uint, todoID uint, title *string, description *string, completed *bool) (*models.Todo, error)
	DeleteTodo(userID uint, todoID uint) error
}

type TodoService struct {
	repo repository.ITodoRepository
}

func NewTodoService(repo repository.ITodoRepository) ITodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) GetTodos(userID uint, limit int, page int) (*repository.PaginatedTodos, error) {
	offset := (page - 1) * limit
	if offset < 0 { offset = 0 }
	return s.repo.GetTodosByUserID(userID, limit, offset)
}

func (s *TodoService) CreateTodo(userID uint, title string, description string) (*models.Todo, error) {
	todo := &models.Todo{UserID: userID, Title: title, Description: description}
	if err := s.repo.CreateTodo(todo); err != nil {
		return nil, domain.ErrInternalServer
	}
	return todo, nil
}

func (s *TodoService) GetTodoByID(id uint, userID uint) (*models.Todo, error) {
	todo, err := s.repo.GetTodoByID(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrTodoNotFound
		}
		return nil, domain.ErrInternalServer
	}
	return todo, nil
}

func (s *TodoService) UpdateTodo(userID uint, todoID uint, title *string, description *string, completed *bool) (*models.Todo, error) {
	todo, err := s.GetTodoByID(todoID, userID)
	if err != nil { return nil, err }

	if title != nil { todo.Title = *title }
	if description != nil { todo.Description = *description }
	if completed != nil { todo.Completed = *completed }

	if err := s.repo.UpdateTodo(todo); err != nil {
		return nil, domain.ErrInternalServer
	}
	return todo, nil
}

func (s *TodoService) DeleteTodo(userID uint, todoID uint) error {
	_, err := s.GetTodoByID(todoID, userID)
	if err != nil { return err }

	if err := s.repo.DeleteTodo(todoID, userID); err != nil {
		return domain.ErrInternalServer
	}
	return nil
}