# Go Project Template

A production-ready Go project template following Clean Architecture and Domain-Driven Design principles, optimized for building scalable applications with PostgreSQL or MySQL.

## Features

- **Modular Domain Architecture** - Domain-based code organization for scalability
- **Clean Architecture** - Separation of concerns with domain, repository, use case, and presentation layers
- **Standardized Error Handling** - Domain errors with proper HTTP status code mapping
- **Dependency Injection Container** - Centralized component wiring with lazy initialization and clean resource management
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
- **Input Validation** - Advanced validation with jellydator/validation library including password strength, email format, and custom rules
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
│   ├── app/                    # Dependency injection container
│   │   ├── di.go
│   │   ├── di_test.go
│   │   └── README.md
│   ├── config/                 # Configuration management
│   │   └── config.go
│   ├── database/               # Database connection and transaction management
│   │   ├── database.go
│   │   └── txmanager.go
│   ├── errors/                 # Standardized domain errors
│   │   └── errors.go
│   ├── http/                   # HTTP server and shared infrastructure
│   │   ├── middleware.go
│   │   ├── response.go
│   │   └── server.go
│   ├── httputil/               # HTTP utility functions
│   │   └── response.go         # JSON responses and error mapping
│   ├── outbox/                 # Outbox domain module
│   │   ├── domain/             # Outbox entities
│   │   │   └── outbox_event.go
│   │   └── repository/         # Outbox data access
│   │       └── outbox_repository.go
│   ├── user/                   # User domain module
│   │   ├── domain/             # User entities and domain errors
│   │   │   └── user.go
│   │   ├── http/               # User HTTP handlers
│   │   │   ├── dto/            # Request/response DTOs
│   │   │   │   ├── request.go
│   │   │   │   ├── response.go
│   │   │   │   └── mapper.go
│   │   │   └── user_handler.go
│   │   ├── repository/         # User data access
│   │   │   └── user_repository.go
│   │   └── usecase/            # User business logic
│   │       └── user_usecase.go
│   ├── validation/             # Custom validation rules
│   │   ├── rules.go
│   │   └── rules_test.go
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

- **`domain/`** - Contains entities, value objects, domain types, and domain-specific errors
- **`usecase/`** - Defines UseCase interfaces and implements business logic and orchestration
- **`repository/`** - Handles data persistence and retrieval, transforms infrastructure errors to domain errors
- **`http/`** - Contains HTTP handlers and DTOs (Data Transfer Objects)
  - **`dto/`** - Request/response DTOs and mappers (API contracts)

### Shared Utilities

- **`app/`** - Dependency injection container for assembling application components
- **`errors/`** - Standardized domain errors for expressing business intent
- **`httputil/`** - Shared HTTP utility functions including error mapping and JSON responses
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
    "password": "SecurePass123!"
  }'
```

**Password Requirements:**
- Minimum 8 characters
- At least one uppercase letter
- At least one lowercase letter
- At least one number
- At least one special character

**Validation Errors:**

If validation fails, you'll receive a 422 Unprocessable Entity response with details:

```json
{
  "error": "invalid_input",
  "message": "email: must be a valid email address; password: password must contain at least one uppercase letter."
}
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

## Error Handling

The project implements a standardized error handling system that expresses business intent rather than exposing infrastructure details.

### Domain Error Architecture

**Standard Domain Errors** (`internal/errors/errors.go`)

The project defines standard domain errors that can be used across all modules:

```go
var (
    ErrNotFound      = errors.New("not found")          // 404 Not Found
    ErrConflict      = errors.New("conflict")           // 409 Conflict
    ErrInvalidInput  = errors.New("invalid input")      // 422 Unprocessable Entity
    ErrUnauthorized  = errors.New("unauthorized")       // 401 Unauthorized
    ErrForbidden     = errors.New("forbidden")          // 403 Forbidden
)
```

**Domain-Specific Errors** (`internal/user/domain/user.go`)

Each domain defines its own specific errors that wrap the standard errors:

```go
var (
    ErrUserNotFound      = errors.Wrap(errors.ErrNotFound, "user not found")
    ErrUserAlreadyExists = errors.Wrap(errors.ErrConflict, "user already exists")
    ErrInvalidEmail      = errors.Wrap(errors.ErrInvalidInput, "invalid email format")
)
```

### Error Flow Through Layers

**1. Repository Layer** - Transforms infrastructure errors to domain errors:

```go
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
    if err := sqlutil.Get(ctx, querier, "users", opts, &user); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, domain.ErrUserNotFound  // Infrastructure → Domain
        }
        return nil, apperrors.Wrap(err, "failed to get user by id")
    }
    return &user, nil
}
```

**2. Use Case Layer** - Returns domain errors directly:

```go
func (uc *UserUseCase) RegisterUser(ctx context.Context, input RegisterUserInput) (*domain.User, error) {
    if strings.TrimSpace(input.Email) == "" {
        return nil, domain.ErrEmailRequired  // Domain error
    }
    
    if err := uc.userRepo.Create(ctx, user); err != nil {
        return nil, err  // Pass through domain errors
    }
    return user, nil
}
```

**3. HTTP Handler Layer** - Maps domain errors to HTTP responses:

```go
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
    user, err := h.userUseCase.RegisterUser(r.Context(), input)
    if err != nil {
        httputil.HandleError(w, err, h.logger)  // Auto-maps to HTTP status
        return
    }
    httputil.MakeJSONResponse(w, http.StatusCreated, response)
}
```

### HTTP Error Responses

The `httputil.HandleError` function automatically maps domain errors to appropriate HTTP status codes:

| Domain Error | HTTP Status | Error Code | Example Response |
|--------------|-------------|------------|------------------|
| `ErrNotFound` | 404 | `not_found` | `{"error":"not_found","message":"The requested resource was not found"}` |
| `ErrConflict` | 409 | `conflict` | `{"error":"conflict","message":"A conflict occurred with existing data"}` |
| `ErrInvalidInput` | 422 | `invalid_input` | `{"error":"invalid_input","message":"invalid email format"}` |
| `ErrUnauthorized` | 401 | `unauthorized` | `{"error":"unauthorized","message":"Authentication is required"}` |
| `ErrForbidden` | 403 | `forbidden` | `{"error":"forbidden","message":"You don't have permission"}` |
| Unknown | 500 | `internal_error` | `{"error":"internal_error","message":"An internal error occurred"}` |

### Benefits

1. **No Infrastructure Leaks** - Database errors are never exposed to API clients
2. **Business Intent** - Errors express domain concepts (`ErrUserNotFound` vs `sql.ErrNoRows`)
3. **Consistent HTTP Mapping** - Same domain error always maps to same HTTP status
4. **Type-Safe** - Use `errors.Is()` to check for specific error types
5. **Structured Responses** - All errors return consistent JSON format
6. **Centralized Logging** - All errors are logged with full context before responding

### Adding Errors to New Domains

When creating a new domain, define domain-specific errors:

```go
// internal/product/domain/product.go
var (
    ErrProductNotFound    = errors.Wrap(errors.ErrNotFound, "product not found")
    ErrInsufficientStock  = errors.Wrap(errors.ErrConflict, "insufficient stock")
    ErrInvalidPrice       = errors.Wrap(errors.ErrInvalidInput, "invalid price")
)
```

Then use `httputil.HandleError()` in your HTTP handlers for automatic mapping.

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

### Dependency Injection Container

The project uses a custom dependency injection (DI) container located in `internal/app/` to manage all application components. This provides:

- **Centralized component wiring** - All dependencies are assembled in one place
- **Lazy initialization** - Components are only created when first accessed
- **Singleton pattern** - Each component is initialized once and reused
- **Clean resource management** - Unified shutdown for all resources
- **Thread-safe** - Safe for concurrent access across goroutines

**Example usage:**

```go
// Create container with configuration
container := app.NewContainer(cfg)

// Get HTTP server (automatically initializes all dependencies)
server, err := container.HTTPServer()
if err != nil {
    return fmt.Errorf("failed to initialize HTTP server: %w", err)
}

// Clean shutdown
defer container.Shutdown(ctx)
```

The container manages the entire dependency graph:

```
Container
├── Infrastructure (Database, Logger)
├── Repositories (User, Outbox)
├── Use Cases (User)
└── Presentation (HTTP Server, Worker)
```

For more details on the DI container, see [`internal/app/README.md`](internal/app/README.md).

### Modular Domain Architecture

The project follows a modular domain-driven structure where each business domain is self-contained:

**User Domain** (`internal/user/`)
- `domain/` - User entity and types
- `usecase/` - UseCase interface and user business logic implementation
- `repository/` - User data persistence
- `http/` - User HTTP endpoints and handlers
  - `dto/` - Request/response DTOs and mappers

**Outbox Domain** (`internal/outbox/`)
- `domain/` - OutboxEvent entity and status types
- `repository/` - Event persistence and retrieval

**Shared Infrastructure**
- `app/` - Dependency injection container for component assembly
- `config/` - Application configuration
- `database/` - Database connection and transaction management
- `errors/` - Standardized domain errors (ErrNotFound, ErrConflict, etc.)
- `http/` - HTTP server, middleware, and shared utilities
- `httputil/` - Reusable HTTP utilities (JSON responses, error handling, status code mapping)
- `worker/` - Background event processing

### Benefits of This Structure

1. **Scalability** - Easy to add new domains without affecting existing code
2. **Encapsulation** - Each domain is self-contained with clear boundaries
3. **Team Collaboration** - Teams can work on different domains independently
4. **Maintainability** - Related code is co-located, making it easier to understand and modify
5. **Dependency Management** - Centralized DI container simplifies component wiring and testing

### Adding New Domains

To add a new domain (e.g., `product`):

#### 1. Create the domain structure

```
internal/product/
├── domain/
│   └── product.go              # Domain entity + domain errors
├── usecase/
│   └── product_usecase.go      # UseCase interface + business logic
├── repository/
│   └── product_repository.go   # Data access (returns domain errors)
└── http/
    ├── dto/
    │   ├── request.go          # API request DTOs
    │   ├── response.go         # API response DTOs
    │   └── mapper.go           # DTO-to-domain mappers
    └── product_handler.go      # HTTP handlers (uses httputil.HandleError)
```

**Define domain errors in your entity file:**

```go
// internal/product/domain/product.go
package domain

import (
    "time"
    apperrors "github.com/yourname/yourproject/internal/errors"
)

type Product struct {
    ID          int64
    Name        string
    Price       float64
    Stock       int
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// Domain-specific errors
var (
    ErrProductNotFound    = apperrors.Wrap(apperrors.ErrNotFound, "product not found")
    ErrInsufficientStock  = apperrors.Wrap(apperrors.ErrConflict, "insufficient stock")
    ErrInvalidPrice       = apperrors.Wrap(apperrors.ErrInvalidInput, "invalid price")
)
```

#### 2. Register in DI container

Add the new domain to the dependency injection container (`internal/app/di.go`):

```go
// Add fields to Container struct
type Container struct {
    // ... existing fields
    productRepo        *productRepository.ProductRepository
    productUseCase     productUsecase.UseCase  // Interface, not concrete type
    productRepoInit    sync.Once
    productUseCaseInit sync.Once
}

// Add getter methods
func (c *Container) ProductRepository() (*productRepository.ProductRepository, error) {
    var err error
    c.productRepoInit.Do(func() {
        c.productRepo, err = c.initProductRepository()
        if err != nil {
            c.initErrors["productRepo"] = err
        }
    })
    // ... error handling
    return c.productRepo, nil
}

func (c *Container) ProductUseCase() (productUsecase.UseCase, error) {
    var err error
    c.productUseCaseInit.Do(func() {
        c.productUseCase, err = c.initProductUseCase()
        if err != nil {
            c.initErrors["productUseCase"] = err
        }
    })
    // ... error handling
    return c.productUseCase, nil
}

// Add initialization methods
func (c *Container) initProductRepository() (*productRepository.ProductRepository, error) {
    db, err := c.DB()
    if err != nil {
        return nil, fmt.Errorf("failed to get database: %w", err)
    }
    return productRepository.NewProductRepository(db, c.config.DBDriver), nil
}

func (c *Container) initProductUseCase() (productUsecase.UseCase, error) {
    txManager, err := c.TxManager()
    if err != nil {
        return nil, fmt.Errorf("failed to get tx manager: %w", err)
    }
    
    productRepo, err := c.ProductRepository()
    if err != nil {
        return nil, fmt.Errorf("failed to get product repository: %w", err)
    }
    
    return productUsecase.NewProductUseCase(txManager, productRepo)
}
```

#### 3. Wire handlers in HTTP server

Update `internal/http/server.go` to register product routes:

```go
productHandler := productHttp.NewProductHandler(container.ProductUseCase(), logger)
mux.HandleFunc("/api/products", productHandler.HandleProducts)
```

**Tips:**
- Define a UseCase interface in your usecase package to enable dependency inversion
- Define domain-specific errors in your domain package by wrapping standard errors
- Repository layer should transform infrastructure errors (like `sql.ErrNoRows`) to domain errors
- Use case layer should return domain errors directly without additional wrapping
- Use `httputil.HandleError()` in HTTP handlers for automatic error-to-HTTP status mapping
- Use the shared `httputil.MakeJSONResponse` function for consistent JSON responses
- Keep domain models free of JSON tags - use DTOs for API serialization
- Implement validation in your request DTOs
- Create mapper functions to convert between DTOs and domain models
- Register all components in the DI container for proper lifecycle management
- HTTP handlers should depend on the UseCase interface, not concrete implementations

### Clean Architecture Layers

1. **Domain Layer** - Contains business entities, domain errors, and rules (e.g., `internal/user/domain`)
2. **Repository Layer** - Data access implementations using sqlutil; transforms infrastructure errors to domain errors (e.g., `internal/user/repository`)
3. **Use Case Layer** - UseCase interfaces and application business logic; returns domain errors (e.g., `internal/user/usecase`)
4. **Presentation Layer** - HTTP handlers that map domain errors to HTTP responses (e.g., `internal/user/http`)
5. **Utility Layer** - Shared utilities including error handling and mapping (e.g., `internal/httputil`, `internal/errors`)

**Dependency Inversion Principle:** The presentation layer (HTTP handlers) and infrastructure (DI container) depend on UseCase interfaces defined in the usecase layer, not on concrete implementations. This enables better testability and decoupling.

**Error Flow:** Errors flow from repository → use case → handler, where they are transformed from infrastructure concerns to domain concepts, and finally to appropriate HTTP responses.

### Data Transfer Objects (DTOs)

The project enforces clear boundaries between internal domain models and external API contracts using DTOs:

**Domain Models** (`internal/user/domain/user.go`)
- Pure internal representation of business entities
- No JSON tags - completely decoupled from API serialization
- Focus on business rules and domain logic

**DTOs** (`internal/user/http/dto/`)
- `request.go` - API request structures with validation
- `response.go` - API response structures
- `mapper.go` - Conversion functions between DTOs and domain models

**Example:**

```go
// Request DTO
type RegisterUserRequest struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

func (r *RegisterUserRequest) Validate() error {
    if r.Name == "" {
        return errors.New("name is required")
    }
    // ... validation logic
}

// Response DTO
type UserResponse struct {
    ID        int64     `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Mapper functions
func ToRegisterUserInput(req RegisterUserRequest) usecase.RegisterUserInput {
    return usecase.RegisterUserInput{
        Name:     req.Name,
        Email:    req.Email,
        Password: req.Password,
    }
}

func ToUserResponse(user *domain.User) UserResponse {
    return UserResponse{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }
}
```

**Benefits:**
1. **Separation of Concerns** - Domain models evolve independently from API contracts
2. **Security** - Sensitive fields (like passwords) are never exposed in API responses
3. **Flexibility** - Different API views of the same domain model (e.g., summary vs detailed)
4. **Versioning** - Easy to maintain multiple API versions with different DTOs
5. **Validation** - Request validation happens at the DTO level before reaching domain logic

### Input Validation

The project uses the [jellydator/validation](https://github.com/jellydator/validation) library for comprehensive input validation at both the DTO and use case layers.

**Custom Validation Rules** (`internal/validation/rules.go`)

The project provides reusable validation rules:

```go
// Password strength validation
PasswordStrength{
    MinLength:      8,
    RequireUpper:   true,
    RequireLower:   true,
    RequireNumber:  true,
    RequireSpecial: true,
}

// Email format validation
Email

// No leading/trailing whitespace
NoWhitespace

// Not blank after trimming
NotBlank
```

**DTO Validation Example:**

```go
func (r *RegisterUserRequest) Validate() error {
    err := validation.ValidateStruct(r,
        validation.Field(&r.Name,
            validation.Required.Error("name is required"),
            appValidation.NotBlank,
            validation.Length(1, 255).Error("name must be between 1 and 255 characters"),
        ),
        validation.Field(&r.Email,
            validation.Required.Error("email is required"),
            appValidation.NotBlank,
            appValidation.Email,
            validation.Length(5, 255).Error("email must be between 5 and 255 characters"),
        ),
        validation.Field(&r.Password,
            validation.Required.Error("password is required"),
            validation.Length(8, 128).Error("password must be between 8 and 128 characters"),
            appValidation.PasswordStrength{
                MinLength:      8,
                RequireUpper:   true,
                RequireLower:   true,
                RequireNumber:  true,
                RequireSpecial: true,
            },
        ),
    )
    return appValidation.WrapValidationError(err)
}
```

**Validation Layers:**

1. **DTO Layer** - Validates API request structure and basic constraints
2. **Use Case Layer** - Validates business logic rules and constraints
3. **Domain Layer** - Defines domain-specific error types

**Error Responses:**

Validation errors are automatically wrapped as `ErrInvalidInput` and return 422 Unprocessable Entity:

```json
{
  "error": "invalid_input",
  "message": "password: password must contain at least one uppercase letter."
}
```

**Benefits:**
- **Declarative** - Validation rules are clear and concise
- **Reusable** - Custom rules can be shared across the application
- **Type-Safe** - Compile-time validation of struct fields
- **Extensible** - Easy to add custom validation rules
- **Consistent** - Same validation logic at DTO and use case layers
- **User-Friendly** - Detailed error messages help API clients fix issues

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
        httputil.HandleError(w, err, h.logger)
        return
    }
    
    httputil.MakeJSONResponse(w, http.StatusOK, product)
}
```

**HandleError** - Automatic domain error to HTTP status code mapping:

```go
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
    user, err := h.userUseCase.RegisterUser(r.Context(), input)
    if err != nil {
        // Automatically maps domain errors to appropriate HTTP status codes
        // ErrNotFound → 404, ErrConflict → 409, ErrInvalidInput → 422, etc.
        httputil.HandleError(w, err, h.logger)
        return
    }
    httputil.MakeJSONResponse(w, http.StatusCreated, response)
}
```

**HandleValidationError** - For JSON decode and validation errors:

```go
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    httputil.HandleValidationError(w, err, h.logger)  // Returns 400 Bad Request
    return
}
```

These utilities ensure consistent response formatting and error handling across all HTTP endpoints.

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
- [validation](https://github.com/jellydator/validation) - Advanced input validation library
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
- github.com/jellydator/validation
- github.com/urfave/cli
- github.com/golang-migrate/migrate