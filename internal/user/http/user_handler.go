// Package http provides HTTP handlers for user-related operations.
package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/allisson/go-project-template/internal/httputil"
	"github.com/allisson/go-project-template/internal/user/domain"
	"github.com/allisson/go-project-template/internal/user/usecase"
)

// UserUseCaseInterface defines the interface for user use case operations
type UserUseCaseInterface interface {
	RegisterUser(ctx context.Context, input usecase.RegisterUserInput) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, id int64) (*domain.User, error)
}

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userUseCase UserUseCaseInterface
	logger      *slog.Logger
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userUseCase UserUseCaseInterface, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		logger:      logger,
	}
}

// RegisterUser handles user registration
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input usecase.RegisterUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		if h.logger != nil {
			h.logger.Error("failed to decode request body", slog.Any("error", err))
		}
		httputil.MakeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	// Validate input
	if input.Name == "" || input.Email == "" || input.Password == "" {
		httputil.MakeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "name, email, and password are required"})
		return
	}

	user, err := h.userUseCase.RegisterUser(r.Context(), input)
	if err != nil {
		if h.logger != nil {
			h.logger.Error("failed to register user", slog.Any("error", err))
		}
		httputil.MakeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to register user"})
		return
	}

	httputil.MakeJSONResponse(w, http.StatusCreated, user)
}
