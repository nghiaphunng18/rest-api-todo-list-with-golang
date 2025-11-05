package repository

import (
	"todo-api/internal/models"

	"gorm.io/gorm"
)

type UserStore struct {
    db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
    return &UserStore{db: db}
}

func (s *UserStore) CreateUser(user *models.User) error {
    return s.db.Create(user).Error
}

func (s *UserStore) GetUserByUsername(username string) (*models.User, error) {
    var user models.User
    err := s.db.Where("username = ?", username).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}
