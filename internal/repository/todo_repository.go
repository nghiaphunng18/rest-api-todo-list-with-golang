package repository

import (
	"todo-api/internal/models"

	"gorm.io/gorm"
)

type PaginatedTodos struct {
	Todos      []models.Todo
	TotalCount int
}

type ITodoRepository interface {
	GetTodosByUserID(userID uint, limit int, offset int) (*PaginatedTodos, error)
	CreateTodo(todo *models.Todo) error
	GetTodoByID(id uint, userID uint) (*models.Todo, error)
	UpdateTodo(todo *models.Todo) error
	DeleteTodo(id uint, userID uint) error
}

type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) ITodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) GetTodosByUserID(userID uint, limit int, offset int) (*PaginatedTodos, error) {
	var todos []models.Todo
	var totalCount int64

	if err := r.db.Model(&models.Todo{}).Where("user_id = ?", userID).Count(&totalCount).Error; err != nil {
		return nil, err
	}

	err := r.db.Where("user_id = ?", userID).
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&todos).Error

	return &PaginatedTodos{Todos: todos, TotalCount: int(totalCount)}, err
}

func (r *TodoRepository) CreateTodo(todo *models.Todo) error {
	return r.db.Create(todo).Error
}

func (r *TodoRepository) GetTodoByID(id uint, userID uint) (*models.Todo, error) {
	var todo models.Todo
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&todo).Error
	return &todo, err
}

func (r *TodoRepository) UpdateTodo(todo *models.Todo) error {
	return r.db.Save(todo).Error
}

func (r *TodoRepository) DeleteTodo(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Todo{}).Error
}
