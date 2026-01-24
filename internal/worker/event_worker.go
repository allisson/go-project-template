// Package worker provides background workers for processing asynchronous tasks.
package worker

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/allisson/go-project-template/internal/database"
	"github.com/allisson/go-project-template/internal/outbox/domain"
)

// Config holds worker configuration
type Config struct {
	Interval      time.Duration
	BatchSize     int
	MaxRetries    int
	RetryInterval time.Duration
}

// OutboxEventRepository interface defines outbox event repository operations
type OutboxEventRepository interface {
	Create(ctx context.Context, event *domain.OutboxEvent) error
	GetPendingEvents(ctx context.Context, limit int) ([]*domain.OutboxEvent, error)
	Update(ctx context.Context, event *domain.OutboxEvent) error
}

// EventWorker processes outbox events
type EventWorker struct {
	config     Config
	txManager  database.TxManager
	outboxRepo OutboxEventRepository
	logger     *slog.Logger
}

// NewEventWorker creates a new EventWorker
func NewEventWorker(
	config Config,
	txManager database.TxManager,
	outboxRepo OutboxEventRepository,
	logger *slog.Logger,
) *EventWorker {
	return &EventWorker{
		config:     config,
		txManager:  txManager,
		outboxRepo: outboxRepo,
		logger:     logger,
	}
}

// Start starts the worker
func (w *EventWorker) Start(ctx context.Context) error {
	if w.logger != nil {
		w.logger.Info("starting event worker",
			slog.Duration("interval", w.config.Interval),
			slog.Int("batch_size", w.config.BatchSize),
		)
	}

	ticker := time.NewTicker(w.config.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			if w.logger != nil {
				w.logger.Info("stopping event worker")
			}
			return ctx.Err()
		case <-ticker.C:
			if err := w.processEvents(ctx); err != nil {
				if w.logger != nil {
					w.logger.Error("failed to process events", slog.Any("error", err))
				}
			}
		}
	}
}

// processEvents retrieves and processes pending events from the outbox in a transaction.
func (w *EventWorker) processEvents(ctx context.Context) error {
	return w.txManager.WithTx(ctx, func(ctx context.Context) error {
		// Get pending events
		events, err := w.outboxRepo.GetPendingEvents(ctx, w.config.BatchSize)
		if err != nil {
			return err
		}

		if len(events) == 0 {
			return nil
		}

		if w.logger != nil {
			w.logger.Info("processing events", slog.Int("count", len(events)))
		}

		for _, event := range events {
			if err := w.processEvent(ctx, event); err != nil {
				if w.logger != nil {
					w.logger.Error("failed to process event",
						slog.Int64("event_id", event.ID),
						slog.String("event_type", event.EventType),
						slog.Any("error", err),
					)
				}

				// Update event as failed
				event.Retries++
				errorMsg := err.Error()
				event.LastError = &errorMsg

				if event.Retries >= w.config.MaxRetries {
					event.Status = domain.OutboxEventStatusFailed
				}

				if err := w.outboxRepo.Update(ctx, event); err != nil {
					return err
				}
				continue
			}

			// Mark event as processed
			now := time.Now()
			event.Status = domain.OutboxEventStatusProcessed
			event.ProcessedAt = &now

			if err := w.outboxRepo.Update(ctx, event); err != nil {
				return err
			}
		}

		return nil
	})
}

// processEvent handles a single outbox event and implements the event processing logic.
func (w *EventWorker) processEvent(ctx context.Context, event *domain.OutboxEvent) error {
	if w.logger != nil {
		w.logger.Info("processing event",
			slog.Int64("event_id", event.ID),
			slog.String("event_type", event.EventType),
		)
	}

	// Parse event payload
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(event.Payload), &payload); err != nil {
		return err
	}

	// Handle different event types
	switch event.EventType {
	case "user.created":
		if w.logger != nil {
			w.logger.Info("user created event",
				slog.Any("payload", payload),
			)
		}
		// In a real application, you might publish this to a message queue,
		// send notifications, update cache, etc.
	default:
		if w.logger != nil {
			w.logger.Warn("unknown event type", slog.String("event_type", event.EventType))
		}
	}

	return nil
}
