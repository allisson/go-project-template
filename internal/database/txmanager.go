// Package database provides database connection management and configuration.
package database

import (
	"context"
	"database/sql"

	"github.com/allisson/sqlutil"
)

// txKey is a context key type for storing database transactions.
type txKey struct{}

// TxManager manages database transactions
type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

// sqlTxManager implements TxManager for SQL databases.
type sqlTxManager struct {
	db *sql.DB
}

// NewTxManager creates a new TxManager
func NewTxManager(db *sql.DB) TxManager {
	return &sqlTxManager{db: db}
}

// WithTx executes the provided function within a database transaction
func (m *sqlTxManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, txKey{}, tx)

	if err := fn(ctx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}

	return tx.Commit()
}

// GetTx retrieves a transaction from context, or falls back to the DB connection
func GetTx(ctx context.Context, db *sql.DB) sqlutil.Querier {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return db
}
