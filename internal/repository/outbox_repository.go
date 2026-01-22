// Package repository provides data persistence implementations for domain entities.
package repository

import (
	"context"
	"database/sql"

	"github.com/allisson/go-project-template/internal/database"
	"github.com/allisson/go-project-template/internal/domain"
	"github.com/allisson/sqlutil"
)

// OutboxEventRepository handles outbox event persistence
type OutboxEventRepository struct {
	db     *sql.DB
	flavor sqlutil.Flavor
}

// NewOutboxEventRepository creates a new OutboxEventRepository
func NewOutboxEventRepository(db *sql.DB, driver string) *OutboxEventRepository {
	flavor := sqlutil.PostgreSQLFlavor
	if driver == "mysql" {
		flavor = sqlutil.MySQLFlavor
	}
	return &OutboxEventRepository{
		db:     db,
		flavor: flavor,
	}
}

// Create inserts a new outbox event
func (r *OutboxEventRepository) Create(ctx context.Context, event *domain.OutboxEvent) error {
	querier := database.GetTx(ctx, r.db)
	return sqlutil.Insert(ctx, querier, r.flavor, "insert", "outbox_events", event)
}

// GetPendingEvents retrieves pending events with limit
func (r *OutboxEventRepository) GetPendingEvents(ctx context.Context, limit int) ([]*domain.OutboxEvent, error) {
	var events []*domain.OutboxEvent
	opts := sqlutil.NewFindAllOptions(r.flavor).
		WithFilter("status", domain.OutboxEventStatusPending).
		WithOrderBy("created_at ASC").
		WithLimit(limit).
		WithForUpdate("SKIP LOCKED")

	querier := database.GetTx(ctx, r.db)
	if err := sqlutil.Select(ctx, querier, "outbox_events", opts, &events); err != nil {
		return nil, err
	}
	return events, nil
}

// Update updates an outbox event
func (r *OutboxEventRepository) Update(ctx context.Context, event *domain.OutboxEvent) error {
	querier := database.GetTx(ctx, r.db)
	return sqlutil.Update(ctx, querier, r.flavor, "update", "outbox_events", event.ID, event)
}
