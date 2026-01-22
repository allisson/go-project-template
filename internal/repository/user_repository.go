// Package repository provides data persistence implementations for domain entities.
package repository

import (
	"context"
	"database/sql"

	"github.com/allisson/go-project-template/internal/database"
	"github.com/allisson/go-project-template/internal/domain"
	"github.com/allisson/sqlutil"
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
	return sqlutil.Insert(ctx, querier, r.flavor, "insert", "users", user)
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User
	opts := sqlutil.NewFindOptions(r.flavor).WithFilter("id", id)
	querier := database.GetTx(ctx, r.db)
	if err := sqlutil.Get(ctx, querier, "users", opts, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	opts := sqlutil.NewFindOptions(r.flavor).WithFilter("email", email)
	querier := database.GetTx(ctx, r.db)
	if err := sqlutil.Get(ctx, querier, "users", opts, &user); err != nil {
		return nil, err
	}
	return &user, nil
}
