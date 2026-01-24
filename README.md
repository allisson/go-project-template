# Go Project Template

A production-ready Go project template following Clean Architecture and Domain-Driven Design principles, optimized for building scalable applications with PostgreSQL or MySQL.

## Features

- **Modular Domain Architecture** - Domain-based code organization for scalability
- **Clean Architecture** - Separation of concerns with domain, repository, use case, and presentation layers
- **Multiple Database Support** - PostgreSQL and MySQL via unified repository layer
- **Database Migrations** - Separate migrations for PostgreSQL and MySQL using golang-migrate
- **Transaction Management** - TxManager interface for handling database transactions
- **Transactional Outbox Pattern** - Event-driven architecture with guaranteed delivery
- **HTTP Server** - Standard library HTTP server with middleware for logging and panic recovery
- **Worker Process** - Background worker for processing outbox events
- **CLI Interface** - urfave/cli for running server, migrations, and worker
- **Health Checks** - Kubernetes-compatible readiness and liveness endpoints
- **Structured Logging** - JSON logs using slog
- **Configuration** - Environment variable based configuration with go-env
- **Password Hashing** - Secure password hashing with Argon2id via go-pwdhash
- **Docker Support** - Multi-stage Dockerfile for minimal container size
- **CI/CD** - GitHub Actions workflow for linting and testing
- **Comprehensive Makefile** - Easy development and deployment commands

## Project Structure

```
go-project-template/
├── cmd/
│   └── app/                    # Application entry point
│       └── main.go
├── internal/
│   ├── config/                 # Configuration management
│   │   └── config.go
│   ├── database/               # Database connection and transaction management
│   │   ├── database.go
│   │   └── txmanager.go
│   ├── http/                   # HTTP server and shared infrastructure
│   │   ├── middleware.go
│   │   ├── response.go
│   │   └── server.go
│   ├── httputil/               # HTTP utility functions
│   │   └── response.go
│   ├── outbox/                 # Outbox domain module
│   │   ├── domain/             # Outbox entities
│   │   │   └── outbox_event.go
│   │   └── repository/         # Outbox data access
│   │       └── outbox_repository.go
│   ├── user/                   # User domain module
│   │   ├── domain/             # User entities
│   │   │   └── user.go
│   │   ├── http/               # User HTTP handlers
│   │   │   └── user_handler.go
│   │   ├── repository/         # User data access
│   │   │   └── user_repository.go
│   │   └── usecase/            # User business logic
│   │       └── user_usecase.go
│   └── worker/                 # Background workers
│       └── event_worker.go
├── migrations/
│   ├── mysql/                  # MySQL migrations
│   └── postgresql/             # PostgreSQL migrations
├── .github/
│   └── workflows/
│       └── ci.yml
├── Dockerfile
├── Makefile
├── go.mod
└── go.sum
```

### Domain Module Structure

The project follows a modular domain architecture where each business domain is organized in its own directory with clear separation of concerns:

- **`domain/`** - Contains entities, value objects, and domain types
- **`usecase/`** - Implements business logic and orchestrates operations
- **`repository/`** - Handles data persistence and retrieval
- **`http/`** - Contains HTTP handlers and request/response types

### Shared Utilities

- **`httputil/`** - Shared HTTP utility functions used across all domain modules (e.g., `MakeJSONResponse`)
- **`config/`** - Application-wide configuration
- **`database/`** - Database connection and transaction management
- **`worker/`** - Background processing infrastructure

This structure makes it easy to add new domains (e.g., `internal/product/`, `internal/order/`) without affecting existing modules.

## Prerequisites

- Go 1.25 or higher
- PostgreSQL 12+ or MySQL 8.0+
- Docker (optional)
- Make (optional, for convenience commands)

## Quick Start

### 1. Clone the repository

```bash
git clone https://github.com/allisson/go-project-template.git
cd go-project-template
```

### 2. Customize the module path

After cloning, you need to update the import paths to match your project:

**Option 1: Using find and sed (Linux/macOS)**
```bash
# Replace with your actual module path
NEW_MODULE="github.com/yourname/yourproject"

# Update go.mod
sed -i "s|github.com/allisson/go-project-template|$NEW_MODULE|g" go.mod

# Update all Go files
find . -type f -name "*.go" -exec sed -i "s|github.com/allisson/go-project-template|$NEW_MODULE|g" {} +
```

**Option 2: Using PowerShell (Windows)**
```powershell
# Replace with your actual module path
$NEW_MODULE = "github.com/yourname/yourproject"

# Update go.mod
(Get-Content go.mod) -replace 'github.com/allisson/go-project-template', $NEW_MODULE | Set-Content go.mod

# Update all Go files
Get-ChildItem -Recurse -Filter *.go | ForEach-Object {
    (Get-Content $_.FullName) -replace 'github.com/allisson/go-project-template', $NEW_MODULE | Set-Content $_.FullName
}
```

**Option 3: Manually**
1. Update the module name in `go.mod`
2. Search and replace `github.com/allisson/go-project-template` with your module path in all `.go` files

After updating, verify the changes and tidy dependencies:
```bash
go mod tidy
```

**Important:** Also update the `.golangci.yml` file to match your new module path:

```yaml
formatters:
  settings:
    goimports:
      local-prefixes:
        - github.com/yourname/yourproject  # Update this line
```

This ensures the linter correctly groups your local imports.

### 3. Install dependencies

```bash
go mod download
```

### 4. Configure environment variables

The application automatically loads environment variables from a `.env` file. Create a `.env` file in your project root (or any parent directory):

```bash
# Database configuration
DB_DRIVER=postgres  # or mysql
DB_CONNECTION_STRING=postgres://user:password@localhost:5432/mydb?sslmode=disable
DB_MAX_OPEN_CONNECTIONS=25
DB_MAX_IDLE_CONNECTIONS=5
DB_CONN_MAX_LIFETIME=5

# Server configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Logging
LOG_LEVEL=info

# Worker configuration
WORKER_INTERVAL=5
WORKER_BATCH_SIZE=10
WORKER_MAX_RETRIES=3
WORKER_RETRY_INTERVAL=1
```

**Note:** The application searches for the `.env` file recursively from the current working directory up to the root directory. This allows you to run the application from any subdirectory and it will still find your `.env` file.

Alternatively, you can export environment variables directly without a `.env` file.

### 5. Start a database (using Docker)

**PostgreSQL:**
```bash
make dev-postgres
```

**MySQL:**
```bash
make dev-mysql
```

### 6. Run database migrations

```bash
make run-migrate
```

### 7. Start the HTTP server

```bash
make run-server
```

The server will be available at http://localhost:8080

### 8. Start the worker (in another terminal)

```bash
make run-worker
```

## Usage

### HTTP Endpoints

#### Health Check
```bash
curl http://localhost:8080/health
```

#### Readiness Check
```bash
curl http://localhost:8080/ready
```

#### Register User
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "securepassword123"
  }'
```

### CLI Commands

The binary supports three commands via urfave/cli:

#### Start HTTP Server
```bash
./bin/app server
```

#### Run Database Migrations
```bash
./bin/app migrate
```

#### Run Event Worker
```bash
./bin/app worker
```

## Development

### Build the application

```bash
make build
```

### Run tests

```bash
make test
```

### Run tests with coverage

```bash
make test-coverage
```

### Run linter

```bash
make lint
```

### Clean build artifacts

```bash
make clean
```

## Docker

### Build Docker image

```bash
make docker-build
```

### Run server in Docker

```bash
make docker-run-server
```

### Run worker in Docker

```bash
make docker-run-worker
```

### Run migrations in Docker

```bash
make docker-run-migrate
```

## Architecture

### Modular Domain Architecture

The project follows a modular domain-driven structure where each business domain is self-contained:

**User Domain** (`internal/user/`)
- `domain/` - User entity and types
- `usecase/` - User registration, authentication logic
- `repository/` - User data persistence
- `http/` - User HTTP endpoints and handlers

**Outbox Domain** (`internal/outbox/`)
- `domain/` - OutboxEvent entity and status types
- `repository/` - Event persistence and retrieval

**Shared Infrastructure**
- `config/` - Application configuration
- `database/` - Database connection and transaction management
- `http/` - HTTP server, middleware, and shared utilities
- `httputil/` - Reusable HTTP utilities (JSON responses, error handling)
- `worker/` - Background event processing

### Benefits of This Structure

1. **Scalability** - Easy to add new domains without affecting existing code
2. **Encapsulation** - Each domain is self-contained with clear boundaries
3. **Team Collaboration** - Teams can work on different domains independently
4. **Maintainability** - Related code is co-located, making it easier to understand and modify

### Adding New Domains

To add a new domain (e.g., `product`):

```
internal/product/
├── domain/
│   └── product.go
├── usecase/
│   └── product_usecase.go
├── repository/
│   └── product_repository.go
└── http/
    └── product_handler.go
```

**Tip:** Use the shared `httputil.MakeJSONResponse` function in your HTTP handlers for consistent JSON responses across all domains.

### Clean Architecture Layers

1. **Domain Layer** - Contains business entities and rules (e.g., `internal/user/domain`)
2. **Repository Layer** - Data access implementations using sqlutil (e.g., `internal/user/repository`)
3. **Use Case Layer** - Application business logic (e.g., `internal/user/usecase`)
4. **Presentation Layer** - HTTP handlers and server (e.g., `internal/user/http`)
5. **Utility Layer** - Shared utilities and helpers (e.g., `internal/httputil`)

### Transaction Management

The template implements a TxManager interface for handling database transactions:

```go
type TxManager interface {
    WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}
```

Transactions are automatically injected into the context and used by repositories.

### HTTP Utilities

The `httputil` package provides shared HTTP utilities used across all domain modules:

**MakeJSONResponse** - Standardized JSON response formatting:

```go
import "github.com/allisson/go-project-template/internal/httputil"

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
    product, err := h.productUseCase.GetProduct(r.Context(), productID)
    if err != nil {
        httputil.MakeJSONResponse(w, http.StatusNotFound, map[string]string{
            "error": "product not found",
        })
        return
    }
    
    httputil.MakeJSONResponse(w, http.StatusOK, product)
}
```

This ensures consistent response formatting across all HTTP endpoints and eliminates code duplication.

### Transactional Outbox Pattern

User registration demonstrates the transactional outbox pattern:

1. User is created in the database
2. `user.created` event is stored in the outbox table (same transaction)
3. Worker picks up pending events and processes them
4. Events are marked as processed or failed

This guarantees that events are never lost and provides at-least-once delivery.

## Configuration

All configuration is done via environment variables. The application automatically loads a `.env` file if present (searching recursively from the current directory up to the root).

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_HOST` | HTTP server host | `0.0.0.0` |
| `SERVER_PORT` | HTTP server port | `8080` |
| `DB_DRIVER` | Database driver (postgres/mysql) | `postgres` |
| `DB_CONNECTION_STRING` | Database connection string | `postgres://user:password@localhost:5432/mydb?sslmode=disable` |
| `DB_MAX_OPEN_CONNECTIONS` | Max open connections | `25` |
| `DB_MAX_IDLE_CONNECTIONS` | Max idle connections | `5` |
| `DB_CONN_MAX_LIFETIME` | Connection max lifetime | `5` |
| `LOG_LEVEL` | Log level (debug/info/warn/error) | `info` |
| `WORKER_INTERVAL` | Worker poll interval | `5` |
| `WORKER_BATCH_SIZE` | Events to process per batch | `10` |
| `WORKER_MAX_RETRIES` | Max retry attempts | `3` |
| `WORKER_RETRY_INTERVAL` | Retry interval | `1` |

## Database Migrations

Migrations are located in `migrations/postgresql` and `migrations/mysql` directories.

### Creating new migrations

1. Create new `.up.sql` and `.down.sql` files with sequential numbering
2. Follow the naming convention: `000003_description.up.sql`

### Running migrations manually

Use the golang-migrate CLI:

```bash
migrate -path migrations/postgresql -database "postgres://user:password@localhost:5432/mydb?sslmode=disable" up
```

## Testing

The project includes a CI workflow that runs tests with PostgreSQL.

### Running tests locally

```bash
go test -v -race ./...
```

### With coverage

```bash
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Dependencies

### Core Libraries

- [go-env](https://github.com/allisson/go-env) - Environment variable configuration
- [godotenv](https://github.com/joho/godotenv) - Loads environment variables from .env files
- [go-pwdhash](https://github.com/allisson/go-pwdhash) - Password hashing with Argon2id
- [sqlutil](https://github.com/allisson/sqlutil) - SQL utilities for unified database access
- [urfave/cli](https://github.com/urfave/cli) - CLI framework
- [golang-migrate](https://github.com/golang-migrate/migrate) - Database migrations

### Database Drivers

- [lib/pq](https://github.com/lib/pq) - PostgreSQL driver
- [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) - MySQL driver

## License

MIT License - see LICENSE file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

This template uses the following excellent Go libraries:
- github.com/allisson/go-env
- github.com/allisson/go-pwdhash
- github.com/allisson/sqlutil
- github.com/urfave/cli
- github.com/golang-migrate/migrate