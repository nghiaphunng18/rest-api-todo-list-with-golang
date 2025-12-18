package services

import (
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