package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	apperrors "github.com/allisson/go-project-template/internal/errors"
	"github.com/allisson/go-project-template/internal/user/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUserRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	tests := []struct {
		name   string
		driver string
	}{
		{
			name:   "create repository with postgres driver",
			driver: "postgres",
		},
		{
			name:   "create repository with mysql driver",
			driver: "mysql",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewUserRepository(db, tt.driver)
			assert.NotNil(t, repo)
			assert.Equal(t, db, repo.db)
		})
	}
}

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	repo := NewUserRepository(db, "postgres")
	ctx := context.Background()

	user := &domain.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashed_password",
	}

	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.Name, user.Email, user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(ctx, user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	repo := NewUserRepository(db, "postgres")
	ctx := context.Background()

	expectedUser := &domain.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  "hashed_password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.Password, expectedUser.CreatedAt, expectedUser.UpdatedAt)

	mock.ExpectQuery("SELECT (.+) FROM users").
		WithArgs(int64(1)).
		WillReturnRows(rows)

	user, err := repo.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	repo := NewUserRepository(db, "postgres")
	ctx := context.Background()

	mock.ExpectQuery("SELECT (.+) FROM users").
		WithArgs(int64(999)).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetByID(ctx, 999)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.True(t, apperrors.Is(err, domain.ErrUserNotFound))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	repo := NewUserRepository(db, "postgres")
	ctx := context.Background()

	expectedUser := &domain.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  "hashed_password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.Password, expectedUser.CreatedAt, expectedUser.UpdatedAt)

	mock.ExpectQuery("SELECT (.+) FROM users").
		WithArgs("john@example.com").
		WillReturnRows(rows)

	user, err := repo.GetByEmail(ctx, "john@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	repo := NewUserRepository(db, "postgres")
	ctx := context.Background()

	mock.ExpectQuery("SELECT (.+) FROM users").
		WithArgs("notfound@example.com").
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetByEmail(ctx, "notfound@example.com")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.True(t, apperrors.Is(err, domain.ErrUserNotFound))
	assert.NoError(t, mock.ExpectationsWereMet())
}
