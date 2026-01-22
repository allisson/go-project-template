package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/allisson/go-project-template/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOutboxEventRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	tests := []struct {
		name   string
		driver string
	}{
		{
			name:   "create repository with postgres driver",
			driver: "postgres",
		},
		{
			name:   "create repository with mysql driver",
			driver: "mysql",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewOutboxEventRepository(db, tt.driver)
			assert.NotNil(t, repo)
			assert.Equal(t, db, repo.db)
		})
	}
}

func TestOutboxEventRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	repo := NewOutboxEventRepository(db, "postgres")
	ctx := context.Background()

	event := &domain.OutboxEvent{
		EventType: "user.created",
		Payload:   `{"id": 1}`,
		Status:    domain.OutboxEventStatusPending,
		Retries:   0,
	}

	mock.ExpectExec("INSERT INTO outbox_events").
		WithArgs(event.EventType, event.Payload, event.Status, event.Retries, event.LastError, event.ProcessedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(ctx, event)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestOutboxEventRepository_GetPendingEvents(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	repo := NewOutboxEventRepository(db, "postgres")
	ctx := context.Background()

	now := time.Now()
	expectedEvents := []*domain.OutboxEvent{
		{
			ID:        1,
			EventType: "user.created",
			Payload:   `{"id": 1}`,
			Status:    domain.OutboxEventStatusPending,
			Retries:   0,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        2,
			EventType: "user.created",
			Payload:   `{"id": 2}`,
			Status:    domain.OutboxEventStatusPending,
			Retries:   0,
			CreatedAt: now.Add(time.Minute),
			UpdatedAt: now.Add(time.Minute),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "event_type", "payload", "status", "retries", "last_error", "processed_at", "created_at", "updated_at"})
	for _, event := range expectedEvents {
		rows.AddRow(event.ID, event.EventType, event.Payload, event.Status, event.Retries, event.LastError, event.ProcessedAt, event.CreatedAt, event.UpdatedAt)
	}

	// Use AnyArg matcher since sqlutil adds multiple query parameters
	mock.ExpectQuery("SELECT (.+) FROM outbox_events").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	events, err := repo.GetPendingEvents(ctx, 10)
	assert.NoError(t, err)
	assert.NotNil(t, events)
	assert.Len(t, events, 2)
	assert.Equal(t, expectedEvents[0].ID, events[0].ID)
	assert.Equal(t, expectedEvents[1].ID, events[1].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestOutboxEventRepository_GetPendingEvents_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	repo := NewOutboxEventRepository(db, "postgres")
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "event_type", "payload", "status", "retries", "last_error", "processed_at", "created_at", "updated_at"})

	mock.ExpectQuery("SELECT (.+) FROM outbox_events").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	events, err := repo.GetPendingEvents(ctx, 10)
	assert.NoError(t, err)
	assert.Len(t, events, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestOutboxEventRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	repo := NewOutboxEventRepository(db, "postgres")
	ctx := context.Background()

	now := time.Now()
	event := &domain.OutboxEvent{
		ID:          1,
		EventType:   "user.created",
		Payload:     `{"id": 1}`,
		Status:      domain.OutboxEventStatusProcessed,
		Retries:     0,
		ProcessedAt: &now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	mock.ExpectExec("UPDATE outbox_events").
		WithArgs(event.EventType, event.Payload, event.Status, event.Retries, event.LastError, event.ProcessedAt, event.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(ctx, event)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestOutboxEventRepository_Update_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	repo := NewOutboxEventRepository(db, "postgres")
	ctx := context.Background()

	event := &domain.OutboxEvent{
		ID:        999,
		EventType: "user.created",
		Payload:   `{"id": 1}`,
		Status:    domain.OutboxEventStatusProcessed,
	}

	updateError := assert.AnError
	mock.ExpectExec("UPDATE outbox_events").
		WithArgs(event.EventType, event.Payload, event.Status, event.Retries, event.LastError, event.ProcessedAt, event.ID).
		WillReturnError(updateError)

	err = repo.Update(ctx, event)
	assert.Error(t, err)
	assert.Equal(t, updateError, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
