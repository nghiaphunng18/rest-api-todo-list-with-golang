package services

import (
	"errors"
	"todo-api/internal/models"
	"todo-api/internal/repository"

	"gorm.io/gorm"
)

var ErrTodoNotFound = errors.New("todo not found")

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

func (s *TodoService) GetTodoByID(id uint, userID uint) (*models.Todo, error) {
    todo, err := s.todoRepo.GetTodoByID(id, userID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrTodoNotFound
        }
        return nil, err
    }
    return todo, nil
}

func (s *TodoService) UpdateTodo(userID uint, todoID uint, title *string, description *string, completed *bool) (*models.Todo, error) {
    // get todo by id
    todo, err := s.GetTodoByID(todoID, userID)
    if err != nil {
        return nil, err
    }

    // update if field is not nil
    if title != nil {
        todo.Title = *title
    }
    if description != nil {
        todo.Description = *description
    }
    if completed != nil {
        todo.Completed = *completed
    }

    if err := s.todoRepo.UpdateTodo(todo); err != nil {
        return nil, err
    }

    return todo, nil
}

func (s *TodoService) DeleteTodo(userID uint, todoID uint) error {
    err := s.todoRepo.DeleteTodo(todoID, userID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ErrTodoNotFound
        }
        return err
    }
    return nil
}
