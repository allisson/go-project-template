package database

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTxManager(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	txManager := NewTxManager(db)
	assert.NotNil(t, txManager)
	assert.IsType(t, &sqlTxManager{}, txManager)
}

func TestWithTx_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	mock.ExpectBegin()
	mock.ExpectCommit()

	txManager := NewTxManager(db)
	ctx := context.Background()

	err = txManager.WithTx(ctx, func(ctx context.Context) error {
		// Verify transaction is in context
		tx := ctx.Value(txKey{})
		assert.NotNil(t, tx)
		assert.IsType(t, &sql.Tx{}, tx)
		return nil
	})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWithTx_RollbackOnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	mock.ExpectBegin()
	mock.ExpectRollback()

	txManager := NewTxManager(db)
	ctx := context.Background()

	testError := assert.AnError
	err = txManager.WithTx(ctx, func(ctx context.Context) error {
		return testError
	})

	assert.Equal(t, testError, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWithTx_BeginError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	beginError := assert.AnError
	mock.ExpectBegin().WillReturnError(beginError)

	txManager := NewTxManager(db)
	ctx := context.Background()

	err = txManager.WithTx(ctx, func(ctx context.Context) error {
		return nil
	})

	assert.Equal(t, beginError, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWithTx_CommitError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	commitError := assert.AnError
	mock.ExpectBegin()
	mock.ExpectCommit().WillReturnError(commitError)

	txManager := NewTxManager(db)
	ctx := context.Background()

	err = txManager.WithTx(ctx, func(ctx context.Context) error {
		return nil
	})

	assert.Equal(t, commitError, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWithTx_RollbackError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	rollbackError := assert.AnError
	mock.ExpectBegin()
	mock.ExpectRollback().WillReturnError(rollbackError)

	txManager := NewTxManager(db)
	ctx := context.Background()

	err = txManager.WithTx(ctx, func(ctx context.Context) error {
		return assert.AnError
	})

	assert.Equal(t, rollbackError, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTx_WithTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	mock.ExpectBegin()
	mock.ExpectCommit()

	txManager := NewTxManager(db)
	ctx := context.Background()

	err = txManager.WithTx(ctx, func(ctx context.Context) error {
		querier := GetTx(ctx, db)
		assert.NotNil(t, querier)
		assert.IsType(t, &sql.Tx{}, querier)
		return nil
	})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTx_WithoutTransaction(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close() //nolint:errcheck

	ctx := context.Background()
	querier := GetTx(ctx, db)

	assert.NotNil(t, querier)
	assert.Equal(t, db, querier)
}
