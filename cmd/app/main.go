// Package main provides the entry point for the application with CLI commands.
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/allisson/go-project-template/internal/config"
	"github.com/allisson/go-project-template/internal/database"
	"github.com/allisson/go-project-template/internal/http"
	outboxRepository "github.com/allisson/go-project-template/internal/outbox/repository"
	userRepository "github.com/allisson/go-project-template/internal/user/repository"
	userUsecase "github.com/allisson/go-project-template/internal/user/usecase"
	"github.com/allisson/go-project-template/internal/worker"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/urfave/cli/v3"
)

// closeDB closes the database connection and logs any errors.
func closeDB(db *sql.DB, logger *slog.Logger) {
	if err := db.Close(); err != nil {
		logger.Error("failed to close the database", slog.Any("error", err))
	}
}

// closeMigrate closes the migration instance and logs any errors.
func closeMigrate(migrate *migrate.Migrate, logger *slog.Logger) {
	sourceError, databaseError := migrate.Close()
	if sourceError != nil || databaseError != nil {
		logger.Error("failed to close the migrate", slog.Any("source_error", sourceError), slog.Any("database_error", databaseError))
	}
}

func main() {
	cmd := &cli.Command{
		Name:    "app",
		Usage:   "Go project template application",
		Version: "1.0.0",
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "Start the HTTP server",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return runServer(ctx)
				},
			},
			{
				Name:  "migrate",
				Usage: "Run database migrations",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return runMigrations()
				},
			},
			{
				Name:  "worker",
				Usage: "Run the event worker",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return runWorker(ctx)
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		slog.Error("application error", slog.Any("error", err))
		os.Exit(1)
	}
}

// runServer starts the HTTP server with graceful shutdown support.
func runServer(ctx context.Context) error {
	// Load configuration
	cfg := config.Load()

	// Setup logger
	logger := setupLogger(cfg.LogLevel)
	logger.Info("starting server", slog.String("version", "1.0.0"))

	// Connect to database
	db, err := database.Connect(database.Config{
		Driver:             cfg.DBDriver,
		ConnectionString:   cfg.DBConnectionString,
		MaxOpenConnections: cfg.DBMaxOpenConnections,
		MaxIdleConnections: cfg.DBMaxIdleConnections,
		ConnMaxLifetime:    cfg.DBConnMaxLifetime,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer closeDB(db, logger)

	// Initialize components
	txManager := database.NewTxManager(db)
	userRepo := userRepository.NewUserRepository(db, cfg.DBDriver)
	outboxRepo := outboxRepository.NewOutboxEventRepository(db, cfg.DBDriver)

	userUseCaseInstance, err := userUsecase.NewUserUseCase(txManager, userRepo, outboxRepo)
	if err != nil {
		return fmt.Errorf("failed to create user use case: %w", err)
	}

	// Create HTTP server
	server := http.NewServer(cfg.ServerHost, cfg.ServerPort, logger, userUseCaseInstance)

	// Setup graceful shutdown
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Start server in goroutine
	serverErr := make(chan error, 1)
	go func() {
		if err := server.Start(ctx); err != nil {
			serverErr <- err
		}
	}()

	// Wait for shutdown signal or server error
	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.DBConnMaxLifetime)
		defer shutdownCancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown failed: %w", err)
		}
	case err := <-serverErr:
		return err
	}

	return nil
}

// runMigrations executes database migrations based on the configured driver.
func runMigrations() error {
	cfg := config.Load()
	logger := setupLogger(cfg.LogLevel)

	logger.Info("running database migrations",
		slog.String("driver", cfg.DBDriver),
	)

	// Determine migration path based on driver
	migrationsPath := "file://migrations/postgresql"
	if cfg.DBDriver == "mysql" {
		migrationsPath = "file://migrations/mysql"
	}

	m, err := migrate.New(migrationsPath, cfg.DBConnectionString)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer closeMigrate(m, logger)

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info("migrations completed successfully")
	return nil
}

// runWorker starts the event worker with graceful shutdown support.
func runWorker(ctx context.Context) error {
	// Load configuration
	cfg := config.Load()

	// Setup logger
	logger := setupLogger(cfg.LogLevel)
	logger.Info("starting worker", slog.String("version", "1.0.0"))

	// Connect to database
	db, err := database.Connect(database.Config{
		Driver:             cfg.DBDriver,
		ConnectionString:   cfg.DBConnectionString,
		MaxOpenConnections: cfg.DBMaxOpenConnections,
		MaxIdleConnections: cfg.DBMaxIdleConnections,
		ConnMaxLifetime:    cfg.DBConnMaxLifetime,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer closeDB(db, logger)

	// Initialize components
	txManager := database.NewTxManager(db)
	outboxRepo := outboxRepository.NewOutboxEventRepository(db, cfg.DBDriver)

	workerConfig := worker.Config{
		Interval:      cfg.WorkerInterval,
		BatchSize:     cfg.WorkerBatchSize,
		MaxRetries:    cfg.WorkerMaxRetries,
		RetryInterval: cfg.WorkerRetryInterval,
	}

	eventWorker := worker.NewEventWorker(workerConfig, txManager, outboxRepo, logger)

	// Setup graceful shutdown
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Start worker
	return eventWorker.Start(ctx)
}

// setupLogger creates and configures a structured logger based on the specified log level.
func setupLogger(level string) *slog.Logger {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})

	return slog.New(handler)
}
