// Package dto provides data transfer objects for the user HTTP layer.
package dto

import "errors"

// RegisterUserRequest represents the API request for user registration
type RegisterUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate validates the RegisterUserRequest
func (r *RegisterUserRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	return nil
}
