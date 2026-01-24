// Package domain defines the core user domain entities and types.
package domain

import (
	"time"

	"github.com/allisson/go-project-template/internal/errors"
)

// User represents a user in the system
type User struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name" fieldtag:"insert,update"`
	Email     string    `db:"email" fieldtag:"insert,update"`
	Password  string    `db:"password" fieldtag:"insert,update"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Domain-specific errors for user operations.
var (
	// ErrUserNotFound indicates the requested user does not exist.
	ErrUserNotFound = errors.Wrap(errors.ErrNotFound, "user not found")

	// ErrUserAlreadyExists indicates a user with the same email already exists.
	ErrUserAlreadyExists = errors.Wrap(errors.ErrConflict, "user already exists")

	// ErrInvalidEmail indicates the email format is invalid.
	ErrInvalidEmail = errors.Wrap(errors.ErrInvalidInput, "invalid email format")

	// ErrInvalidPassword indicates the password doesn't meet requirements.
	ErrInvalidPassword = errors.Wrap(errors.ErrInvalidInput, "invalid password")

	// ErrNameRequired indicates the name field is required.
	ErrNameRequired = errors.Wrap(errors.ErrInvalidInput, "name is required")

	// ErrEmailRequired indicates the email field is required.
	ErrEmailRequired = errors.Wrap(errors.ErrInvalidInput, "email is required")

	// ErrPasswordRequired indicates the password field is required.
	ErrPasswordRequired = errors.Wrap(errors.ErrInvalidInput, "password is required")
)
