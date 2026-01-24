package worker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/allisson/go-project-template/internal/outbox/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTxManager is a mock implementation of database.TxManager
type MockTxManager struct {
	mock.Mock
}

func (m *MockTxManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	// Execute the function to test the logic inside
	return fn(ctx)
}

// MockOutboxEventRepository is a mock implementation of repository.OutboxEventRepository
type MockOutboxEventRepository struct {
	mock.Mock
}

func (m *MockOutboxEventRepository) Create(ctx context.Context, event *domain.OutboxEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockOutboxEventRepository) GetPendingEvents(ctx context.Context, limit int) ([]*domain.OutboxEvent, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.OutboxEvent), args.Error(1)
}

func (m *MockOutboxEventRepository) Update(ctx context.Context, event *domain.OutboxEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func TestNewEventWorker(t *testing.T) {
	config := Config{
		Interval:      5 * time.Second,
		BatchSize:     10,
		MaxRetries:    3,
		RetryInterval: 1 * time.Minute,
	}
	txManager := &MockTxManager{}
	outboxRepo := &MockOutboxEventRepository{}

	worker := NewEventWorker(config, txManager, outboxRepo, nil)

	assert.NotNil(t, worker)
	assert.Equal(t, config.Interval, worker.config.Interval)
	assert.Equal(t, config.BatchSize, worker.config.BatchSize)
	assert.Equal(t, config.MaxRetries, worker.config.MaxRetries)
}

func TestEventWorker_Start_ContextCancellation(t *testing.T) {
	config := Config{
		Interval:      100 * time.Millisecond,
		BatchSize:     10,
		MaxRetries:    3,
		RetryInterval: 1 * time.Minute,
	}
	txManager := &MockTxManager{}
	outboxRepo := &MockOutboxEventRepository{}

	worker := NewEventWorker(config, txManager, outboxRepo, nil)

	ctx, cancel := context.WithCancel(context.Background())

	// Cancel context immediately
	cancel()

	err := worker.Start(ctx)
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
}

func TestEventWorker_ProcessEvents_Success(t *testing.T) {
	config := Config{
		Interval:      5 * time.Second,
		BatchSize:     10,
		MaxRetries:    3,
		RetryInterval: 1 * time.Minute,
	}
	txManager := &MockTxManager{}
	outboxRepo := &MockOutboxEventRepository{}

	worker := NewEventWorker(config, txManager, outboxRepo, nil)

	ctx := context.Background()
	uuid1 := uuid.Must(uuid.NewV7())
	uuid2 := uuid.Must(uuid.NewV7())
	events := []*domain.OutboxEvent{
		{
			ID:        uuid1,
			EventType: "user.created",
			Payload:   `{"user_id": 1, "name": "John Doe", "email": "john@example.com"}`,
			Status:    domain.OutboxEventStatusPending,
			Retries:   0,
		},
		{
			ID:        uuid2,
			EventType: "user.created",
			Payload:   `{"user_id": 2, "name": "Jane Doe", "email": "jane@example.com"}`,
			Status:    domain.OutboxEventStatusPending,
			Retries:   0,
		},
	}

	// Setup expectations
	txManager.On("WithTx", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	outboxRepo.On("GetPendingEvents", ctx, config.BatchSize).Return(events, nil)
	outboxRepo.On("Update", ctx, mock.MatchedBy(func(e *domain.OutboxEvent) bool {
		return e.Status == domain.OutboxEventStatusProcessed && e.ProcessedAt != nil
	})).Return(nil).Times(2)

	err := worker.processEvents(ctx)

	assert.NoError(t, err)
	txManager.AssertExpectations(t)
	outboxRepo.AssertExpectations(t)
}

func TestEventWorker_ProcessEvents_NoEvents(t *testing.T) {
	config := Config{
		Interval:      5 * time.Second,
		BatchSize:     10,
		MaxRetries:    3,
		RetryInterval: 1 * time.Minute,
	}
	txManager := &MockTxManager{}
	outboxRepo := &MockOutboxEventRepository{}

	worker := NewEventWorker(config, txManager, outboxRepo, nil)

	ctx := context.Background()
	emptyEvents := []*domain.OutboxEvent{}

	// Setup expectations
	txManager.On("WithTx", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	outboxRepo.On("GetPendingEvents", ctx, config.BatchSize).Return(emptyEvents, nil)

	err := worker.processEvents(ctx)

	assert.NoError(t, err)
	txManager.AssertExpectations(t)
	outboxRepo.AssertExpectations(t)
}

func TestEventWorker_ProcessEvents_GetPendingError(t *testing.T) {
	config := Config{
		Interval:      5 * time.Second,
		BatchSize:     10,
		MaxRetries:    3,
		RetryInterval: 1 * time.Minute,
	}
	txManager := &MockTxManager{}
	outboxRepo := &MockOutboxEventRepository{}

	worker := NewEventWorker(config, txManager, outboxRepo, nil)

	ctx := context.Background()
	getError := errors.New("database error")

	// Setup expectations
	txManager.On("WithTx", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	outboxRepo.On("GetPendingEvents", ctx, config.BatchSize).Return(nil, getError)

	err := worker.processEvents(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	txManager.AssertExpectations(t)
	outboxRepo.AssertExpectations(t)
}

func TestEventWorker_ProcessEvents_InvalidJSON(t *testing.T) {
	config := Config{
		Interval:      5 * time.Second,
		BatchSize:     10,
		MaxRetries:    3,
		RetryInterval: 1 * time.Minute,
	}
	txManager := &MockTxManager{}
	outboxRepo := &MockOutboxEventRepository{}

	worker := NewEventWorker(config, txManager, outboxRepo, nil)

	ctx := context.Background()
	uuid1 := uuid.Must(uuid.NewV7())
	events := []*domain.OutboxEvent{
		{
			ID:        uuid1,
			EventType: "user.created",
			Payload:   `invalid json`,
			Status:    domain.OutboxEventStatusPending,
			Retries:   0,
		},
	}

	// Setup expectations
	txManager.On("WithTx", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	outboxRepo.On("GetPendingEvents", ctx, config.BatchSize).Return(events, nil)
	outboxRepo.On("Update", ctx, mock.MatchedBy(func(e *domain.OutboxEvent) bool {
		return e.ID == uuid1 && e.Retries == 1 && e.LastError != nil
	})).Return(nil)

	err := worker.processEvents(ctx)

	assert.NoError(t, err) // processEvents should not return error, just log and update event
	txManager.AssertExpectations(t)
	outboxRepo.AssertExpectations(t)
}

func TestEventWorker_ProcessEvents_MaxRetriesReached(t *testing.T) {
	config := Config{
		Interval:      5 * time.Second,
		BatchSize:     10,
		MaxRetries:    3,
		RetryInterval: 1 * time.Minute,
	}
	txManager := &MockTxManager{}
	outboxRepo := &MockOutboxEventRepository{}

	worker := NewEventWorker(config, txManager, outboxRepo, nil)

	ctx := context.Background()
	uuid1 := uuid.Must(uuid.NewV7())
	events := []*domain.OutboxEvent{
		{
			ID:        uuid1,
			EventType: "user.created",
			Payload:   `invalid json`,
			Status:    domain.OutboxEventStatusPending,
			Retries:   2, // Will become 3 after this attempt
		},
	}

	// Setup expectations
	txManager.On("WithTx", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	outboxRepo.On("GetPendingEvents", ctx, config.BatchSize).Return(events, nil)
	outboxRepo.On("Update", ctx, mock.MatchedBy(func(e *domain.OutboxEvent) bool {
		return e.ID == uuid1 &&
			e.Retries == 3 &&
			e.Status == domain.OutboxEventStatusFailed &&
			e.LastError != nil
	})).Return(nil)

	err := worker.processEvents(ctx)

	assert.NoError(t, err)
	txManager.AssertExpectations(t)
	outboxRepo.AssertExpectations(t)
}

func TestEventWorker_ProcessEvents_UpdateError(t *testing.T) {
	config := Config{
		Interval:      5 * time.Second,
		BatchSize:     10,
		MaxRetries:    3,
		RetryInterval: 1 * time.Minute,
	}
	txManager := &MockTxManager{}
	outboxRepo := &MockOutboxEventRepository{}

	worker := NewEventWorker(config, txManager, outboxRepo, nil)

	ctx := context.Background()
	uuid1 := uuid.Must(uuid.NewV7())
	events := []*domain.OutboxEvent{
		{
			ID:        uuid1,
			EventType: "user.created",
			Payload:   `{"user_id": 1}`,
			Status:    domain.OutboxEventStatusPending,
			Retries:   0,
		},
	}

	updateError := errors.New("update failed")

	// Setup expectations
	txManager.On("WithTx", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	outboxRepo.On("GetPendingEvents", ctx, config.BatchSize).Return(events, nil)
	outboxRepo.On("Update", ctx, mock.AnythingOfType("*domain.OutboxEvent")).Return(updateError)

	err := worker.processEvents(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update failed")
	txManager.AssertExpectations(t)
	outboxRepo.AssertExpectations(t)
}

func TestEventWorker_ProcessEvent_Success(t *testing.T) {
	config := Config{
		Interval:      5 * time.Second,
		BatchSize:     10,
		MaxRetries:    3,
		RetryInterval: 1 * time.Minute,
	}
	txManager := &MockTxManager{}
	outboxRepo := &MockOutboxEventRepository{}

	worker := NewEventWorker(config, txManager, outboxRepo, nil)

	ctx := context.Background()
	uuid1 := uuid.Must(uuid.NewV7())
	event := &domain.OutboxEvent{
		ID:        uuid1,
		EventType: "user.created",
		Payload:   `{"user_id": 1, "name": "John Doe", "email": "john@example.com"}`,
		Status:    domain.OutboxEventStatusPending,
		Retries:   0,
	}

	err := worker.processEvent(ctx, event)

	assert.NoError(t, err)
}

func TestEventWorker_ProcessEvent_UnknownEventType(t *testing.T) {
	config := Config{
		Interval:      5 * time.Second,
		BatchSize:     10,
		MaxRetries:    3,
		RetryInterval: 1 * time.Minute,
	}
	txManager := &MockTxManager{}
	outboxRepo := &MockOutboxEventRepository{}

	worker := NewEventWorker(config, txManager, outboxRepo, nil)

	ctx := context.Background()
	uuid1 := uuid.Must(uuid.NewV7())
	event := &domain.OutboxEvent{
		ID:        uuid1,
		EventType: "unknown.event",
		Payload:   `{"data": "test"}`,
		Status:    domain.OutboxEventStatusPending,
		Retries:   0,
	}

	err := worker.processEvent(ctx, event)

	assert.NoError(t, err) // Unknown events are just logged as warning
}

func TestEventWorker_ProcessEvent_InvalidJSON(t *testing.T) {
	config := Config{
		Interval:      5 * time.Second,
		BatchSize:     10,
		MaxRetries:    3,
		RetryInterval: 1 * time.Minute,
	}
	txManager := &MockTxManager{}
	outboxRepo := &MockOutboxEventRepository{}

	worker := NewEventWorker(config, txManager, outboxRepo, nil)

	ctx := context.Background()
	uuid1 := uuid.Must(uuid.NewV7())
	event := &domain.OutboxEvent{
		ID:        uuid1,
		EventType: "user.created",
		Payload:   `invalid json`,
		Status:    domain.OutboxEventStatusPending,
		Retries:   0,
	}

	err := worker.processEvent(ctx, event)

	assert.Error(t, err)
}
