package services

import (
	"todo-api/internal/models"
	"todo-api/internal/repository"
)

type TodoService struct {
    todoRepo *repository.TodoRepository
}

func NewTodoService(todoRepo *repository.TodoRepository) *TodoService {
    return &TodoService{todoRepo: todoRepo}
}

func (s *TodoService) GetTodos(userID uint, limit int, page int) (*repository.PaginatedTodos, error) {
    // page starts from 1
	offset := (page - 1) * limit
    
    if offset < 0 {
        offset = 0
    }
    
    return s.todoRepo.GetTodosByUserID(userID, limit, offset)
}

func (s *TodoService) CreateTodo(userID uint, title string, description string) (*models.Todo, error) {
    todo := &models.Todo{
        UserID:      userID,
        Title:       title,
        Description: description,
        Completed:   false,
    }

    if err := s.todoRepo.CreateTodo(todo); err != nil {
        return nil, err
    }

    return todo, nil
}