// Package repository provides data persistence implementations for user entities.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/allisson/go-project-template/internal/database"
	"github.com/allisson/go-project-template/internal/user/domain"
	"github.com/allisson/sqlutil"

	apperrors "github.com/allisson/go-project-template/internal/errors"
)

// UserRepository handles user persistence
type UserRepository struct {
	db     *sql.DB
	flavor sqlutil.Flavor
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sql.DB, driver string) *UserRepository {
	flavor := sqlutil.PostgreSQLFlavor
	if driver == "mysql" {
		flavor = sqlutil.MySQLFlavor
	}
	return &UserRepository{
		db:     db,
		flavor: flavor,
	}
}

// Create inserts a new user
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	querier := database.GetTx(ctx, r.db)
	if err := sqlutil.Insert(ctx, querier, r.flavor, "insert", "users", user); err != nil {
		// Check for unique constraint violation (duplicate email)
		if isUniqueViolation(err) {
			return domain.ErrUserAlreadyExists
		}
		return apperrors.Wrap(err, "failed to create user")
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User
	opts := sqlutil.NewFindOptions(r.flavor).WithFilter("id", id)
	querier := database.GetTx(ctx, r.db)
	if err := sqlutil.Get(ctx, querier, "users", opts, &user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, apperrors.Wrap(err, "failed to get user by id")
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	opts := sqlutil.NewFindOptions(r.flavor).WithFilter("email", email)
	querier := database.GetTx(ctx, r.db)
	if err := sqlutil.Get(ctx, querier, "users", opts, &user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, apperrors.Wrap(err, "failed to get user by email")
	}
	return &user, nil
}

// isUniqueViolation checks if the error is a unique constraint violation.
// This works for both PostgreSQL and MySQL drivers.
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	errMsg := strings.ToLower(err.Error())
	// PostgreSQL: "duplicate key value violates unique constraint" or "pq: duplicate key"
	// MySQL: "Error 1062: Duplicate entry"
	return strings.Contains(errMsg, "duplicate key") ||
		strings.Contains(errMsg, "unique constraint") ||
		strings.Contains(errMsg, "duplicate entry") ||
		strings.Contains(errMsg, "1062")
}
