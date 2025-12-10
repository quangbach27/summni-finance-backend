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
- **BankDetails**: Value object encapsulating bank-specific information

### Detailed Domain Model Analysis

#### Wallet Entity (Aggregate Root)

```go
type Wallet struct {
    id           ID  // wallet.ID which is uuid.UUID
    name         string
    isStrictMode bool
    currency     valueobject.Currency
    allocations  []*Allocation
}
```

**Key Features:**

- **Aggregate Root**: Manages `Allocation` entities with consistency boundaries
- **Rich Domain Behavior**: `TopUp()`, `Withdraw()`, `TotalBalance()` methods
- **Currency Consistency**: Enforces same currency across all operations
- **Strict Mode**: Business rule for operation validation
- **Factory Pattern**: `NewWallet(name, currency, isStrictMode, allocations)` with validation, `UnmarshalWalletFromDB()` for persistence reconstruction

**Business Rules Enforced:**

- Wallet must have a name and at least one allocation ("wallet is not belong to assert source")
- Currency consistency across all operations with detailed error messages
- Asset source must exist in wallet allocations for top-up/withdraw operations
- Negative money amounts prevented at Money creation and operation level
- Immutable ID generation using UUIDv7 for both Wallet and AssetSource
- Strict validation for empty names and nil IDs

#### AssetSource Entity

```go
type AssetSource struct {
    id          ID  // assetsource.ID which is uuid.UUID
    name        string
    sourceType  SourceType  // TypeBank or TypeCash (struct with code field)
    balance     valueobject.Money
    ownerID     uuid.UUID
    bankDetails *BankDetails // Nil for cash assets
}
```

**Key Features:**

- **Type Safety**: Distinct factories for Bank vs Cash assets
- **SourceType**: Struct with code field, predefined constants `TypeBank = SourceType{code: "BANK"}` and `TypeCash = SourceType{code: "CASH"}`
- **Conditional Composition**: Bank details only for bank assets
- **Value Object Integration**: Uses `BankDetails` value object for bank-specific data
- **Owner Association**: Links to user via `ownerID`

**Type-Safe Factories:**

```go
func NewBankAssetSource(ownerID uuid.UUID, name string, initBalance valueobject.Money, bankName, accountNumber string) (*AssetSource, error)
func NewCashAssetSource(ownerID uuid.UUID, name string, initBalance valueobject.Money) (*AssetSource, error)
```

#### Allocation Entity (Part of Wallet Aggregate)

```go
type Allocation struct {
    assetSourceID assetsource.ID
    amount        valueobject.Money
}
```

**Key Features:**

- **Entity within Aggregate**: Managed by Wallet aggregate root
- **Linking Mechanism**: Connects wallets to asset sources
- **Value Object Usage**: Money with currency validation
- **Factory Method**: `NewAllocation(assetSourceID, amount)` with validation

## **Design Patterns**

### 1. CQRS (Command Query Responsibility Segregation)

```go
// Commands for state changes
type CreateAssetSourceHandler cqrs.CommandHandler[CreateAssetSourceCmd]
type createAssetSourceHandler struct{}

// Queries for data retrieval
type GetAssetSourceHandler cqrs.QueryHandler[GetAssetSourceCmd, AssetSource]
type getAssetSourceHandler struct{}
```

**Current Implementation:**

- **CreateAssetSourceCmd**: Empty struct, handler logs creation event
- **GetAssetSourceCmd**: Empty struct, returns hardcoded AssetSource example
- **Application Layer**: Aggregates Commands and Queries structs
- **Auto-applied Decorators**: `ApplyCommandDecorators()` and `ApplyQueryDecorator()` for logging
- **Factory Pattern**: `NewCreateAssetSourceHandler()` and `NewGetAssetSoureHandler()` return decorated handlers

**Usage Guidelines:**

- Commands: Use for operations that modify state
- Queries: Use for read-only operations
- Apply logging decorators for observability
- Return types defined in query/types.go for DTOs

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
func NewBankAssetSource(ownerID uuid.UUID, name string, initBalance valueobject.Money, bankName, accountNumber string) (*AssetSource, error)
func NewCashAssetSource(ownerID uuid.UUID, name string, initBalance valueobject.Money) (*AssetSource, error)
func NewAllocation(assetSourceID assetsource.ID, amount valueobject.Money) (*Allocation, error)
func NewMoney(amount int64, currency Currency) (Money, error)
func NewBankDetails(bankName, accountNumber string) (BankDetails, error)
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
func (w *Wallet) TopUp(assetSourceID assetsource.ID, amount valueobject.Money) error
func (w *Wallet) Withdraw(assetSourceID assetsource.ID, amount valueobject.Money) error
func (w *Wallet) TotalBalance() (valueobject.Money, error)
func (m Money) Add(other Money) (Money, error)
func (m Money) Subtract(other Money) (Money, error)
func (m Money) LessOrEqualThan(other Money) bool
func (m Money) IsZero() bool
```

**Domain-specific Validations:**

- **TopUp**: Validates positive amounts and currency consistency, finds existing allocation
- **Withdraw**: Checks sufficient funds and currency matching
- **Money Operations**: Prevent negative results and enforce currency consistency
- **Allocation Factory**: `NewAllocation()` validates non-nil asset source ID and positive amount

**BankDetails Value Object:**

```go
type BankDetails struct {
    bankName      string
    accountNumber string
}

func NewBankDetails(bankName, accountNumber string) (BankDetails, error)
```

- Immutable value object for bank-specific information
- Validation ensures both bank name and account number are required
- Used only in bank-type asset sources (nil for cash assets)

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

### Current Implementation Status & Next Steps

**✅ Completed Features:**

- Complete domain model with Wallet, AssetSource, and Allocation entities
- Rich value objects (Money, Currency, BankDetails) with proper validation
- CQRS pattern foundation with command/query handlers
- Factory patterns for type-safe object creation
- Clean Architecture with proper layer separation
- HTTP server setup with middleware stack
- Configuration management with environment variables
- Structured logging with request context
- Docker development environment with hot reload

**🚧 In Progress / Needs Implementation:**

- **Repository Pattern**: Interfaces defined but implementations needed
- **Database Integration**: Migrations ready but repository layer incomplete
- **CQRS Handlers**: Currently return hardcoded responses, need actual business logic
- **Persistence Layer**: Domain entities need serialization/deserialization
- **API Endpoints**: HTTP handlers need actual business logic integration
- **Transaction Support**: Multi-operation consistency not yet implemented
- **Input Validation**: HTTP request validation layer missing
- **Error Handling**: Domain errors need proper HTTP status mapping

**📋 Recommended Implementation Order:**

1. Complete repository interfaces and implementations
2. Implement proper CQRS command/query handlers with actual business logic
3. Add database persistence layer with proper domain reconstruction
4. Implement HTTP request/response DTOs and validation
5. Add transaction management for multi-aggregate operations
6. Implement comprehensive error handling with proper HTTP responses
7. Add authentication and authorization middleware
8. Implement event sourcing for audit trails

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

## **SOLID Principles Analysis**

### Assessment: **Well-Designed, Not Over-Complicated**

#### ✅ **Single Responsibility Principle (SRP)**

- **Wallet**: Manages financial allocations and operations
- **AssetSource**: Represents funding sources (bank/cash)
- **Money**: Handles monetary calculations with currency validation
- **BankDetails**: Encapsulates bank-specific information
- **Each class has one reason to change**

#### ✅ **Open/Closed Principle (OCP)**

- **AssetSource**: Extensible via factory pattern (Bank/Cash types)
- **CQRS Handlers**: New commands/queries can be added without modifying existing ones
- **Middleware**: Decorator pattern allows adding new behaviors

#### ✅ **Liskov Substitution Principle (LSP)**

- **Interface-based design**: All implementations can substitute their interfaces
- **Value Objects**: Immutable and behaviorally consistent
- **Factory methods**: Return same interface contracts

#### ✅ **Interface Segregation Principle (ISP)**

- **Ports package**: Specific interfaces (FinanceServerInterface)
- **CQRS interfaces**: Separate CommandHandler and QueryHandler
- **No forced implementation of unused methods**

#### ✅ **Dependency Inversion Principle (DIP)**

- **Clean Architecture**: Domain doesn't depend on infrastructure
- **Dependency Injection**: Constructor injection throughout
- **Interfaces in ports**: Abstract dependencies properly

### Design Complexity Assessment: **Appropriate**

**Pros:**

- **Domain Complexity**: Financial applications require strict business rules
- **Type Safety**: Strong typing prevents runtime errors
- **Maintainability**: Clear separation enables easy testing and modification
- **Scalability**: Architecture supports growth without major refactoring

**Could Be Simplified For:**

- **Simple CRUD applications** (but this isn't one)
- **Prototype/MVP stages** (but production-ready code requires rigor)

**Conclusion**: The complexity level is **justified** for a financial domain where data integrity, business rule enforcement, and maintainability are critical.

## **SQL Security Considerations**

### Current Status: **Repository Pattern Not Yet Implemented**

**Security Guidelines for Future Database Integration:**

#### 1. **SQL Injection Prevention**

```go
// ✅ Use parameterized queries
func (r *walletRepository) GetByID(ctx context.Context, id WalletID) (*Wallet, error) {
    query := `SELECT id, name, currency, is_strict_mode FROM wallets WHERE id = $1`
    row := r.db.QueryRowContext(ctx, query, id)
    // ...
}

// ❌ NEVER use string concatenation
func (r *walletRepository) GetByName(name string) {
    query := fmt.Sprintf("SELECT * FROM wallets WHERE name = '%s'", name) // VULNERABLE
}
```

#### 2. **Input Validation at Domain Layer**

```go
// ✅ Already implemented - validate at domain boundaries
func NewWallet(name string, currency valueobject.Currency, ...) (*Wallet, error) {
    if name == "" {
        return nil, errors.New("wallet name cannot be empty")
    }
    // Prevent SQL injection through domain validation
}
```

#### 3. **Database Access Control**

```go
// ✅ Recommended patterns for future implementation
type WalletRepository interface {
    GetByOwner(ctx context.Context, ownerID uuid.UUID) ([]*Wallet, error)    // Owner-scoped queries
    Create(ctx context.Context, wallet *Wallet) error                        // Input validation
    UpdateBalance(ctx context.Context, walletID WalletID, amount Money) error // Atomic operations
}
```

#### 4. **Connection Security**

```go
// ✅ Configuration already prepared
type DatabaseConfig struct {
    host     string  // Use TLS connections
    database string  // Principle of least privilege DB names
    user     string  // Dedicated service user (not admin)
    password string  // Strong passwords, prefer environment variables
}
```

#### 5. **Query Security Patterns**

- **Parameterized Queries**: Use `$1, $2` placeholders
- **Owner Scoping**: Always include `owner_id` in financial queries
- **Transaction Boundaries**: Use database transactions for multi-table operations
- **Audit Logging**: Log all financial operations with user context
- **Rate Limiting**: Prevent abuse of expensive queries

#### 6. **Recommended Libraries**

```go
import (
    "github.com/jmoiron/sqlx"          // Enhanced database/sql
    "github.com/lib/pq"                // PostgreSQL driver
    "github.com/golang-migrate/migrate" // Schema migrations
)
```

#### 7. **Security Headers & Middleware** (Already Implemented)

```go
// ✅ Current security measures
router.Use(
    middleware.SetHeader("X-Content-Type-Options", "nosniff"),
    middleware.SetHeader("X-Frame-Options", "deny"),
)
```

### Migration Security Checklist

- [ ] Use migrations with version control
- [ ] Review all SQL scripts for injection vulnerabilities
- [ ] Implement row-level security (RLS) for multi-tenant data
- [ ] Set up database connection pooling with proper timeouts
- [ ] Configure SSL/TLS for database connections
- [ ] Implement audit trails for financial transactions
