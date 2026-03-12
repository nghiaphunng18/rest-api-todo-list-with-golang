package repository_test

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"todo-api/internal/models"
	"todo-api/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open sqlmock: %s", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	if err != nil {
		t.Fatalf("Failed to open gorm: %s", err)
	}

	return gormDB, mock
}

func TestTodoRepository_CreateTodo(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repository.NewTodoRepository(db)

	t.Run("Create_Success", func(t *testing.T) {
		todo := &models.Todo{
			UserID: 1,
			Title:  "Test SQL",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `todos`")).
			WithArgs(todo.UserID, todo.Title, "", false, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.CreateTodo(todo)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestTodoRepository_GetTodoByID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repository.NewTodoRepository(db)

	t.Run("GetByID_Found", func(t *testing.T) {
		todoID := uint(1)
		userID := uint(1)

		rows := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "completed", "created_at", "updated_at"}).
			AddRow(todoID, userID, "Mocked Todo", "", false, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `todos` WHERE id = ? AND user_id = ? ORDER BY `todos`.`id` LIMIT ?")).
			WithArgs(todoID, userID, 1).
			WillReturnRows(rows)

		result, err := repo.GetTodoByID(todoID, userID)

		assert.NoError(t, err)
		assert.Equal(t, "Mocked Todo", result.Title)
		assert.Equal(t, todoID, result.ID)
	})

	t.Run("GetByID_NotFound", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `todos` WHERE id = ? AND user_id = ?")).
			WithArgs(99, 1, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		_, err := repo.GetTodoByID(99, 1)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	})
}