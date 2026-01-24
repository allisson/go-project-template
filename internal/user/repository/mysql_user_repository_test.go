package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apperrors "github.com/allisson/go-project-template/internal/errors"
	"github.com/allisson/go-project-template/internal/user/domain"
)

func TestNewMySQLUserRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	repo := NewMySQLUserRepository(db)
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

func TestMySQLUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	repo := NewMySQLUserRepository(db)
	ctx := context.Background()

	uuid1 := uuid.Must(uuid.NewV7())
	user := &domain.User{
		ID:       uuid1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashed_password",
	}

	uuidBytes, _ := uuid1.MarshalBinary()
	mock.ExpectExec("INSERT INTO users").
		WithArgs(uuidBytes, user.Name, user.Email, user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(ctx, user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLUserRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	repo := NewMySQLUserRepository(db)
	ctx := context.Background()

	uuid1 := uuid.Must(uuid.NewV7())
	expectedUser := &domain.User{
		ID:        uuid1,
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  "hashed_password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	uuidBytes, _ := uuid1.MarshalBinary()
	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
		AddRow(uuidBytes, expectedUser.Name, expectedUser.Email, expectedUser.Password, expectedUser.CreatedAt, expectedUser.UpdatedAt)

	mock.ExpectQuery("SELECT (.+) FROM users").
		WithArgs(uuidBytes).
		WillReturnRows(rows)

	user, err := repo.GetByID(ctx, uuid1)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLUserRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	repo := NewMySQLUserRepository(db)
	ctx := context.Background()

	notFoundUUID := uuid.Must(uuid.NewV7())
	uuidBytes, _ := notFoundUUID.MarshalBinary()
	mock.ExpectQuery("SELECT (.+) FROM users").
		WithArgs(uuidBytes).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetByID(ctx, notFoundUUID)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.True(t, apperrors.Is(err, domain.ErrUserNotFound))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLUserRepository_GetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	repo := NewMySQLUserRepository(db)
	ctx := context.Background()

	uuid1 := uuid.Must(uuid.NewV7())
	expectedUser := &domain.User{
		ID:        uuid1,
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  "hashed_password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	uuidBytes, _ := uuid1.MarshalBinary()
	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
		AddRow(uuidBytes, expectedUser.Name, expectedUser.Email, expectedUser.Password, expectedUser.CreatedAt, expectedUser.UpdatedAt)

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

func TestMySQLUserRepository_GetByEmail_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	repo := NewMySQLUserRepository(db)
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
