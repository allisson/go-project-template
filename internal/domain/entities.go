// Package domain defines the core domain entities and types for the application.
package domain

import "time"

// User represents a user in the system
type User struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name" fieldtag:"insert,update"`
	Email     string    `db:"email" json:"email" fieldtag:"insert,update"`
	Password  string    `db:"password" json:"-" fieldtag:"insert,update"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// OutboxEventStatus represents the status of an outbox event
type OutboxEventStatus string

const (
	OutboxEventStatusPending   OutboxEventStatus = "pending"
	OutboxEventStatusProcessed OutboxEventStatus = "processed"
	OutboxEventStatusFailed    OutboxEventStatus = "failed"
)

// OutboxEvent represents an event in the transactional outbox pattern
type OutboxEvent struct {
	ID          int64             `db:"id" json:"id"`
	EventType   string            `db:"event_type" json:"event_type" fieldtag:"insert,update"`
	Payload     string            `db:"payload" json:"payload" fieldtag:"insert,update"`
	Status      OutboxEventStatus `db:"status" json:"status" fieldtag:"insert,update"`
	Retries     int               `db:"retries" json:"retries" fieldtag:"insert,update"`
	LastError   *string           `db:"last_error" json:"last_error,omitempty" fieldtag:"insert,update"`
	ProcessedAt *time.Time        `db:"processed_at" json:"processed_at,omitempty" fieldtag:"insert,update"`
	CreatedAt   time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time         `db:"updated_at" json:"updated_at"`
}
