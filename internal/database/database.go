package database

import (
	"fmt"
	"todo-api/internal/config"
	"todo-api/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	// connect to db
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	
	err = db.AutoMigrate(&models.User{}, &models.Todo{}) 
    if err != nil {
        return nil, fmt.Errorf("failed to migrate database: %w", err)
    }
	
	return db, nil
}
