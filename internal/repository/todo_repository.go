package repository

import (
	"todo-api/internal/models"

	"gorm.io/gorm"
)

type PaginatedTodos struct {
    Todos      []models.Todo
    TotalCount int 
}

type TodoRepository struct {
    db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
    return &TodoRepository{db: db}
}

func (r *TodoRepository) GetTodosByUserID(userID uint, limit int, offset int) (*PaginatedTodos, error) {
    var todos []models.Todo
    var totalCount int64 

    if err := r.db.Model(&models.Todo{}).
                   Where("user_id = ?", userID).
                   Count(&totalCount).Error; err != nil {
        return nil, err
    }

    
    err := r.db.Where("user_id = ?", userID).
              Limit(limit).
              Offset(offset).
              Order("created_at DESC").
              Find(&todos).Error

    if err != nil {
        return nil, err
    }

    return &PaginatedTodos{
        Todos:      todos,
        TotalCount: int(totalCount), 
    }, nil
}

func (r *TodoRepository) CreateTodo(todo *models.Todo) error {
    return r.db.Create(todo).Error
}

func (r *TodoRepository) GetTodoByID(id uint, userID uint) (*models.Todo, error) {
    var todo models.Todo
    if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&todo).Error; err != nil {
        return nil, err
    }
    return &todo, nil
}

func (r *TodoRepository) UpdateTodo(todo *models.Todo) error {
    return r.db.Save(todo).Error
}
