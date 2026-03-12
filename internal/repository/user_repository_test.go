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

func setupUserMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

func TestUserRepository_CreateUser(t *testing.T) {
	db, mock := setupUserMockDB(t)
	repo := repository.NewUserRepository(db)

	t.Run("CreateUser_Success", func(t *testing.T) {
		user := &models.User{
			Username:     "nghia_test",
			PasswordHash: "hashed_password",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
			WithArgs(user.Username, user.PasswordHash, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.CreateUser(user)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_GetUserByUsername(t *testing.T) {
	db, mock := setupUserMockDB(t)
	repo := repository.NewUserRepository(db)

	t.Run("GetByUsername_Found", func(t *testing.T) {
		username := "test_user"
		
		rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "created_at", "updated_at"}).
			AddRow(1, username, "hashed_pass", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ? ORDER BY `users`.`id` LIMIT ?")).
			WithArgs(username, 1).
			WillReturnRows(rows)

		result, err := repo.GetUserByUsername(username)

		assert.NoError(t, err)
		assert.Equal(t, username, result.Username)
		assert.Equal(t, uint(1), result.ID)
	})

	t.Run("GetByUsername_NotFound", func(t *testing.T) {
		username := "unknown"

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ?")).
			WithArgs(username, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		_, err := repo.GetUserByUsername(username)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	})
}