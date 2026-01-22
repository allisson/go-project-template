// Package http provides HTTP server implementation and request handlers.
package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/allisson/go-project-template/internal/usecase"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userUseCase *usecase.UserUseCase
	logger      *slog.Logger
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userUseCase *usecase.UserUseCase, logger *slog.Logger) *UserHandler {
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
		h.logger.Error("failed to decode request body", slog.Any("error", err))
		makeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	user, err := h.userUseCase.RegisterUser(r.Context(), input)
	if err != nil {
		h.logger.Error("failed to register user", slog.Any("error", err))
		makeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to register user"})
		return
	}

	makeJSONResponse(w, http.StatusCreated, user)
}
