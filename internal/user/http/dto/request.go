// Package dto provides data transfer objects for the user HTTP layer.
package dto

import "github.com/allisson/go-project-template/internal/user/domain"

// RegisterUserRequest represents the API request for user registration
type RegisterUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate validates the RegisterUserRequest
// Note: This provides basic JSON structure validation.
// Detailed validation is handled by the use case layer.
func (r *RegisterUserRequest) Validate() error {
	if r.Name == "" {
		return domain.ErrNameRequired
	}
	if r.Email == "" {
		return domain.ErrEmailRequired
	}
	if r.Password == "" {
		return domain.ErrPasswordRequired
	}
	return nil
}
