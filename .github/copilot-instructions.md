# Sumni Finance Backend - Code Structure & Design Patterns

This document outlines the architecture, design patterns, and technologies used in the Sumni Finance Backend application to provide better context for code suggestions and development.

## **Overall Architecture: Clean Architecture / Hexagonal Architecture**

The codebase follows **Clean Architecture** principles with clear separation of concerns:

```
├── cmd/server/           # Application entry point
├── internal/
│   ├── common/          # Shared cross-cutting concerns
│   │   ├── cqrs/       # CQRS pattern implementation
│   │   ├── logs/       # Structured logging
│   │   ├── server/     # HTTP server setup
│   │   └── valueobject/ # Domain value objects
│   ├── config/         # Application configuration
│   └── finance/        # Finance domain module
│       ├── app/        # Application layer (CQRS handlers)
│       ├── domain/     # Domain entities & business logic
│       ├── ports/      # Adapters/Interfaces
│       └── service/    # Service layer
```

## **Domain-Driven Design (DDD)**

### Bounded Context

- **Finance Domain**: Main business domain handling financial operations
- **Common Domain**: Shared utilities and cross-cutting concerns

### Domain Entities

- **Wallet**: Aggregate root managing financial allocations
- **AssetSource**: Bank accounts, cash, and other funding sources
- **Transaction**: Financial operations (IN/OUT/TRANSFER)
- **Allocation**: Links between wallets and asset sources

### Value Objects

- **Money**: Immutable monetary values with currency validation
- **Currency**: Supported currencies (USD, VND, KRW)
- **IDs**: Type-safe identifiers using UUID

## **Design Patterns**

### 1. CQRS (Command Query Responsibility Segregation)

```go
// Commands for state changes
type CreateAssetSourceHandler cqrs.CommandHandler[CreateAssetSourceCmd]

// Queries for data retrieval
type GetAssetSourceHandler cqrs.QueryHandler[GetAssetSourceCmd, AssetSource]
```

**Usage Guidelines:**

- Commands: Use for operations that modify state
- Queries: Use for read-only operations
- Apply logging decorators for observability

### 2. Decorator Pattern

```go
// Automatic logging for all CQRS operations
func ApplyCommandDecorators[C any](handler CommandHandler[C]) CommandHandler[C]
func ApplyQueryDecorator[Q any, R any](handler QueryHandler[Q, R]) QueryHandler[Q, R]
```

**Applied in:**

- CQRS logging decorators
- HTTP middleware stack

### 3. Factory Pattern

```go
func NewBankAssetSource(id int64, name string, balance valueobject.Money, ownerID uuid.UUID, accountNumber, bankName string) (*BankAssetSource, error)
func NewCashAssetSource() (*CashAssetSource, error)
func NewMoney(amount int64, currency Currency) (Money, error)
```

**Usage Guidelines:**

- Always return error for validation failures
- Enforce business invariants at creation time

### 4. Repository Pattern (Interface-based)

- Interfaces defined in `ports/` package
- Implementations in infrastructure layer (planned)

### 5. Dependency Injection

```go
// Constructor injection throughout
func NewFinanceServerInterface(app app.Application) ports.FinanceServerInterface
func NewApplication() app.Application
```

### 6. Singleton Pattern

```go
// Configuration singleton with thread-safety
func GetConfig() *Config {
    once.Do(func() {
        cfg := loadConfig()
        configInstance = cfg
    })
    return configInstance
}
```

## **Technology Stack**

### Core Technologies

- **Language**: Go 1.24.3
- **HTTP Router**: Chi v5 - Fast, idiomatic HTTP router
- **UUID Generation**: Google UUID v6 (using UUIDv7 for time-ordered IDs)
- **Testing**: Testify for assertions and test suites

### HTTP & Web Layer

- **Middleware**: Chi middleware stack with custom logging
- **CORS**: go-chi/cors for cross-origin requests
- **Rendering**: go-chi/render for JSON responses
- **Logging**: Go's built-in `log/slog` package with structured logging

### Development & Deployment

- **Containerization**: Docker with Docker Compose
- **Hot Reload**: Reflex for development with file watching
- **Debugging**: Delve (dlv) debugger with remote debugging support
- **Linting**: golangci-lint with comprehensive rule set
- **Build Tools**: Make for task automation

## **Key Architectural Decisions**

### 1. Error Handling Strategy

```go
// Custom error types with HTTP context
type SlugError struct {
    error     string
    slug      string
    errorType ErrorType
}

// Error type categories
var (
    ErrorTypeUnknown        = ErrorType{"unknown"}
    ErrorTypeAuthorization  = ErrorType{"authorization"}
    ErrorTypeIncorrectInput = ErrorType{"incorrect-input"}
)
```

**Guidelines:**

- Use `SlugError` for user-facing errors
- Include error slugs for client-side error handling
- Categorize errors by type for proper HTTP status codes

### 2. Configuration Management

```go
type Config struct {
    database DatabaseConfig
    app      AppConfig
}
```

**Features:**

- Environment-based configuration with defaults
- Singleton pattern with lazy initialization
- Type-safe configuration access with getters

### 3. Structured Logging Strategy

```go
// Request-scoped logging with context
func GetLogEntry(r *http.Request) *slog.Logger
func LoggerFromCtx(ctx context.Context) *slog.Logger
```

**Guidelines:**

- Use request-scoped loggers for traceability
- Include request IDs in all log entries
- Apply consistent log levels (DEBUG, INFO, WARN, ERROR)

### 4. HTTP Layer Architecture

```go
// Graceful server setup with timeouts
server := &http.Server{
    Addr:              addr,
    ReadHeaderTimeout: 500 * time.Millisecond,
    ReadTimeout:       500 * time.Millisecond,
    IdleTimeout:       time.Second,
    Handler:           rootRouter,
}
```

**Features:**

- Graceful shutdown handling
- Security headers and CORS configuration
- Request timeout and rate limiting
- Comprehensive middleware stack

### 5. Domain Modeling Principles

```go
// Rich domain objects with behavior
func (w *Wallet) TopUp(assetSourceID AssetSourceID, amount valueobject.Money) error
func (w *Wallet) Withdraw(assetSourceID AssetSourceID, amount valueobject.Money) error
func (m Money) Add(other Money) (Money, error)
func (m Money) Subtract(other Money) (Money, error)
```

**Guidelines:**

- Encapsulate business logic in domain entities
- Enforce invariants through method implementations
- Use value objects for primitive obsession prevention
- Return errors for business rule violations

## **Coding Standards & Conventions**

### 1. Package Organization

```go
// Package naming: domain-specific, lowercase
package finance   // ✓ Good
package Finance   // ✗ Avoid

// Import organization: standard, external, internal
import (
    "context"
    "errors"

    "github.com/google/uuid"

    "sumni-finance-backend/internal/common/valueobject"
)
```

### 2. Interface Design

```go
// Interfaces in consumer packages
type FinanceServerInterface interface {
    CreateAssetSource(w http.ResponseWriter, r *http.Request)
    GetAssetSources(w http.ResponseWriter, r *http.Request)
}
```

### 3. Error Handling Patterns

```go
// Always handle errors explicitly
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// Use business-specific errors in domain layer
if amount.Currency() != w.currency {
    return fmt.Errorf("wallet currency is %s but amount is %s", w.currency.Code(), amount.Currency().Code())
}
```

### 4. Testing Conventions

```go
// Test file naming: *_test.go
func TestGetConfig_Singleton(t *testing.T) {
    // Use testify assertions
    assert.Equal(t, expected, actual)
}
```

## **Development Workflow**

### Available Commands

```bash
make dev          # Start development server with hot reload
make dev DEBUG=true # Start with debugger attached
make test         # Run all tests
make stop         # Stop development containers
make lint         # Run code linting
```

### Development Environment Features

- **Hot Reload**: Reflex watches for Go file changes and rebuilds
- **Debug Support**: Delve debugger integration on port 40000
- **Docker Development**: Complete containerized development environment
- **Environment Variables**: Configuration via .env files
- **Testing**: Comprehensive test setup with environment isolation

## **Future Considerations**

### Planned Enhancements

- Database integration (PostgreSQL configuration ready)
- Repository implementations for data persistence
- Event sourcing for transaction history
- API documentation with OpenAPI/Swagger
- Metrics and monitoring integration
- Authentication and authorization middleware

### Scalability Considerations

- Microservice boundaries clearly defined by domain
- Event-driven architecture preparation with CQRS
- Database connection pooling and optimization
- Caching layer integration points identified
- Horizontal scaling support through stateless design

---

This architecture provides a solid foundation for a financial application with proper separation of concerns, maintainable code structure, and excellent development experience. When adding new features, follow the established patterns and maintain the clean architecture principles.
