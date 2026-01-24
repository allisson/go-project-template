package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/allisson/go-project-template/internal/outbox/domain"
)

func TestNewMySQLOutboxEventRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	repo := NewMySQLOutboxEventRepository(db)
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

func TestMySQLOutboxEventRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	repo := NewMySQLOutboxEventRepository(db)
	ctx := context.Background()

	uuid1 := uuid.Must(uuid.NewV7())
	event := &domain.OutboxEvent{
		ID:        uuid1,
		EventType: "user.created",
		Payload:   `{"id": 1}`,
		Status:    domain.OutboxEventStatusPending,
		Retries:   0,
	}

	idBytes, _ := uuid1.MarshalBinary()
	mock.ExpectExec("INSERT INTO outbox_events").
		WithArgs(idBytes, event.EventType, event.Payload, event.Status, event.Retries, event.LastError, event.ProcessedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(ctx, event)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLOutboxEventRepository_GetPendingEvents(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	repo := NewMySQLOutboxEventRepository(db)
	ctx := context.Background()

	now := time.Now()
	uuid1 := uuid.Must(uuid.NewV7())
	uuid2 := uuid.Must(uuid.NewV7())

	idBytes1, _ := uuid1.MarshalBinary()
	idBytes2, _ := uuid2.MarshalBinary()

	rows := sqlmock.NewRows([]string{"id", "event_type", "payload", "status", "retries", "last_error", "processed_at", "created_at", "updated_at"}).
		AddRow(idBytes1, "user.created", `{"id": 1}`, domain.OutboxEventStatusPending, 0, nil, nil, now, now).
		AddRow(idBytes2, "user.created", `{"id": 2}`, domain.OutboxEventStatusPending, 0, nil, nil, now.Add(time.Minute), now.Add(time.Minute))

	mock.ExpectQuery("SELECT (.+) FROM outbox_events").
		WithArgs(domain.OutboxEventStatusPending, 10).
		WillReturnRows(rows)

	events, err := repo.GetPendingEvents(ctx, 10)
	assert.NoError(t, err)
	assert.NotNil(t, events)
	assert.Len(t, events, 2)
	assert.Equal(t, uuid1, events[0].ID)
	assert.Equal(t, uuid2, events[1].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLOutboxEventRepository_GetPendingEvents_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	repo := NewMySQLOutboxEventRepository(db)
	ctx := context.Background()

	rows := sqlmock.NewRows(
		[]string{
			"id",
			"event_type",
			"payload",
			"status",
			"retries",
			"last_error",
			"processed_at",
			"created_at",
			"updated_at",
		},
	)

	mock.ExpectQuery("SELECT (.+) FROM outbox_events").
		WithArgs(domain.OutboxEventStatusPending, 10).
		WillReturnRows(rows)

	events, err := repo.GetPendingEvents(ctx, 10)
	assert.NoError(t, err)
	assert.Len(t, events, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLOutboxEventRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	repo := NewMySQLOutboxEventRepository(db)
	ctx := context.Background()

	now := time.Now()
	uuid1 := uuid.Must(uuid.NewV7())
	event := &domain.OutboxEvent{
		ID:          uuid1,
		EventType:   "user.created",
		Payload:     `{"id": 1}`,
		Status:      domain.OutboxEventStatusProcessed,
		Retries:     0,
		ProcessedAt: &now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	idBytes, _ := uuid1.MarshalBinary()
	mock.ExpectExec("UPDATE outbox_events").
		WithArgs(event.EventType, event.Payload, event.Status, event.Retries, event.LastError, event.ProcessedAt, idBytes).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(ctx, event)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLOutboxEventRepository_Update_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	repo := NewMySQLOutboxEventRepository(db)
	ctx := context.Background()

	uuid1 := uuid.Must(uuid.NewV7())
	event := &domain.OutboxEvent{
		ID:        uuid1,
		EventType: "user.created",
		Payload:   `{"id": 1}`,
		Status:    domain.OutboxEventStatusProcessed,
	}

	idBytes, _ := uuid1.MarshalBinary()
	updateError := assert.AnError
	mock.ExpectExec("UPDATE outbox_events").
		WithArgs(event.EventType, event.Payload, event.Status, event.Retries, event.LastError, event.ProcessedAt, idBytes).
		WillReturnError(updateError)

	err = repo.Update(ctx, event)
	assert.Error(t, err)
	assert.Equal(t, updateError, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
