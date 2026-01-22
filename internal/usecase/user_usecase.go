// Package usecase implements the application's business logic and orchestrates domain operations.
package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/allisson/go-project-template/internal/database"
	"github.com/allisson/go-project-template/internal/domain"
	"github.com/allisson/go-project-template/internal/repository"
	"github.com/allisson/go-pwdhash"
)

// RegisterUserInput contains the input data for user registration
type RegisterUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserUseCase handles user-related business logic
type UserUseCase struct {
	txManager      database.TxManager
	userRepo       *repository.UserRepository
	outboxRepo     *repository.OutboxEventRepository
	passwordHasher *pwdhash.PasswordHasher
}

// NewUserUseCase creates a new UserUseCase
func NewUserUseCase(
	txManager database.TxManager,
	userRepo *repository.UserRepository,
	outboxRepo *repository.OutboxEventRepository,
) (*UserUseCase, error) {
	// Initialize password hasher with interactive policy for user passwords
	hasher, err := pwdhash.New(pwdhash.WithPolicy(pwdhash.PolicyInteractive))
	if err != nil {
		return nil, fmt.Errorf("failed to create password hasher: %w", err)
	}

	return &UserUseCase{
		txManager:      txManager,
		userRepo:       userRepo,
		outboxRepo:     outboxRepo,
		passwordHasher: hasher,
	}, nil
}

// RegisterUser registers a new user and creates a user.created event
func (uc *UserUseCase) RegisterUser(ctx context.Context, input RegisterUserInput) (*domain.User, error) {
	// Hash the password
	hashedPassword, err := uc.passwordHasher.Hash([]byte(input.Password))
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
	}

	// Execute within a transaction
	err = uc.txManager.WithTx(ctx, func(ctx context.Context) error {
		// Create user
		if err := uc.userRepo.Create(ctx, user); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		// Create user.created event payload
		eventPayload := map[string]interface{}{
			"user_id": user.ID,
			"name":    user.Name,
			"email":   user.Email,
		}
		payloadJSON, err := json.Marshal(eventPayload)
		if err != nil {
			return fmt.Errorf("failed to marshal event payload: %w", err)
		}

		// Create outbox event
		outboxEvent := &domain.OutboxEvent{
			EventType: "user.created",
			Payload:   string(payloadJSON),
			Status:    domain.OutboxEventStatusPending,
			Retries:   0,
		}

		if err := uc.outboxRepo.Create(ctx, outboxEvent); err != nil {
			return fmt.Errorf("failed to create outbox event: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (uc *UserUseCase) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return uc.userRepo.GetByEmail(ctx, email)
}

// GetUserByID retrieves a user by ID
func (uc *UserUseCase) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	return uc.userRepo.GetByID(ctx, id)
}
