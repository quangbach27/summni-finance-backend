# Copilot Instructions for Sumni Finance Backend

This document provides context and guidelines for GitHub Copilot when assisting with code development in this project.

---

## 🏗️ Project Architecture

This project follows **Clean Architecture** with **Domain-Driven Design (DDD)** and **Hexagonal Architecture** (Ports & Adapters).

### Architecture Layers

```
┌─────────────────────────────────────────────────────────────┐
│                        HTTP/External                         │
│                     (Framework & Drivers)                    │
└───────────────────────────┬─────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────┐
│                          Ports                               │
│              (Interface Adapters - HTTP Handlers)            │
└───────────────────────────┬─────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────┐
│                      Application Layer                       │
│            (Use Cases - Commands & Queries)                  │
│                      CQRS Pattern                            │
└───────────────────────────┬─────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────┐
│                       Domain Layer                           │
│          (Entities, Value Objects, Repositories)             │
│              Pure Business Logic - No Dependencies           │
└───────────────────────────┬─────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────┐
│                       Adapters                               │
│         (Infrastructure - Database, External APIs)           │
└─────────────────────────────────────────────────────────────┘
```

### Domain Structure

1. **Finance Domain** (`internal/finance/`)

   - Core business logic for financial operations
   - Entities: AssetSource, Wallet, Transaction
   - Value Objects: Money, Currency
   - CQRS: Commands for writes, Queries for reads

2. **Auth Domain** (`internal/auth/`)

   - Keycloak integration (OAuth2/OIDC)
   - Token management and verification
   - Authentication middleware

3. **Common Domain** (`internal/common/`)
   - Shared utilities (CQRS, DB, Logging, Validation)
   - Cross-cutting concerns
   - Reusable value objects

---

## 🎯 Code Generation Guidelines

### When generating new features:

#### 1. Start with Domain Layer

```go
// internal/finance/domain/{entity}/
// - Define entity struct with business rules
// - Create value objects (immutable, validated)
// - Define repository interface
// - Write unit tests
```

#### 2. Application Layer (CQRS)

```go
// Commands (Write Operations)
// internal/finance/app/command/
type CreateEntityHandler struct {
    repo domain.EntityRepository
}

func (h CreateEntityHandler) Handle(ctx context.Context, cmd CreateEntity) error {
    // Validate, create domain entity, persist
}

// Queries (Read Operations)
// internal/finance/app/query/
type GetEntityHandler struct {
    queries *store.Queries
}

func (h GetEntityHandler) Handle(ctx context.Context, q GetEntity) (*EntityDTO, error) {
    // Fetch and return data
}
```

#### 3. Adapter Layer

```go
// internal/finance/adapter/db/
// - Implement repository interface
// - Use SQLC for type-safe queries
// - SQL queries in adapter/db/store/queries/
```

#### 4. Ports Layer

```go
// internal/finance/ports/
// - HTTP handlers
// - Request/Response DTOs
// - Route registration
```

---

## 📝 Code Style & Conventions

### Naming Conventions

- **Packages**: lowercase, no underscores (e.g., `assetsource`, not `asset_source`)
- **Interfaces**: `-er` suffix when appropriate (Reader, Handler)
- **Repository methods**: Domain language (e.g., `SaveAssetSource`, not `InsertAssetSource`)
- **Handlers**: `{Action}{Entity}Handler` (e.g., `CreateAssetSourceHandler`)

### Error Handling

```go
// ✅ Always wrap errors with context
if err != nil {
    return fmt.Errorf("failed to create asset source: %w", err)
}

// ✅ Custom domain errors
var ErrAssetSourceNotFound = errors.New("asset source not found")

// ✅ Error types for specific cases
type ValidationError struct {
    Field   string
    Message string
}
```

### Context Usage

```go
// ✅ Always pass context as first parameter
func (h Handler) Handle(ctx context.Context, cmd Command) error

// ✅ Propagate context to all downstream calls
entity, err := h.repo.Find(ctx, id)
```

### Value Objects

```go
// ✅ Immutable, validated in constructor
func NewMoney(amount decimal.Decimal, currency Currency) (Money, error) {
    if amount.IsNegative() {
        return Money{}, errors.New("amount cannot be negative")
    }
    return Money{amount: amount, currency: currency}, nil
}

// ❌ No setters on value objects
```

### Repository Pattern

```go
// ✅ Interface in domain layer
type AssetSourceRepository interface {
    Save(ctx context.Context, source *AssetSource) error
    FindByID(ctx context.Context, id uuid.UUID) (*AssetSource, error)
}

// ✅ Implementation in adapter layer
type assetsourceRepo struct {
    pool    *pgxpool.Pool
    queries *store.Queries
}
```

---

## 🗄️ Database & SQLC

### SQLC Usage

- All queries in `internal/finance/adapter/db/store/queries/*.sql`
- Run `sqlc generate` after adding/modifying queries
- Use generated code in repository implementations

```sql
-- name: GetAssetSource :one
SELECT * FROM asset_sources WHERE id = $1;

-- name: CreateAssetSource :one
INSERT INTO asset_sources (id, name, type, details)
VALUES ($1, $2, $3, $4)
RETURNING *;
```

### Migrations

- Location: `db/migrations/`
- Format: `YYYYMMDDHHMMSS_description.up.sql` / `*.down.sql`
- Always provide both up and down migrations

---

## 🧪 Testing Guidelines

### Test Structure

```go
func TestCreateAssetSource(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateInput
        want    *AssetSource
        wantErr bool
    }{
        {
            name: "valid input creates asset source",
            // ...
        },
        {
            name: "invalid input returns error",
            // ...
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Test Coverage

- Unit tests for domain entities and value objects
- Integration tests for repositories
- Handler tests with mocked dependencies

---

## 🔒 Security

### Authentication

- Keycloak for authentication (OAuth2/OIDC)
- Token verification via middleware
- Protected routes use `authHandler.AuthMiddleware`

### Best Practices

- Never log sensitive data (tokens, passwords)
- Use parameterized queries (SQLC handles this)
- Validate all external inputs
- No hardcoded secrets (use environment variables)

---

## 🚀 Common Tasks

### Adding a New Entity

1. **Domain**: Create entity in `internal/finance/domain/{entity}/`

   ```go
   type Entity struct {
       id   uuid.UUID
       name string
       // fields...
   }

   // Factory function with validation
   func NewEntity(...) (*Entity, error)
   ```

2. **Repository Interface**: In domain directory

   ```go
   type EntityRepository interface {
       Save(ctx context.Context, e *Entity) error
       FindByID(ctx context.Context, id uuid.UUID) (*Entity, error)
   }
   ```

3. **Migration**: Create in `db/migrations/`

4. **SQLC Queries**: Add to `adapter/db/store/queries/`

5. **Repository Implementation**: In `adapter/db/`

6. **Commands/Queries**: In `app/command/` and `app/query/`

7. **HTTP Handlers**: In `ports/`

8. **Wire Up**: Update `service/application.go` and `ports/http.go`

### Adding a New Endpoint

1. Define handler in `internal/finance/ports/`
2. Add route in `ports/http.go`
3. Update `cmd/server/main.go` if needed

---

## 🛠️ Development Workflow

### Local Development

```bash
make dev              # Start with hot reload
make dev DEBUG=true   # Start with debugger
make test             # Run tests
make stop             # Stop containers
```

### Before Committing

```bash
go fmt ./...
go vet ./...
golangci-lint run
go test -race ./...
```

---

## 📚 Key Principles

1. **Dependency Rule**: Dependencies point inward (domain has no external deps)
2. **CQRS**: Separate models for reads (queries) and writes (commands)
3. **Repository Pattern**: Abstract data access, interface in domain
4. **Value Objects**: Immutable, self-validating
5. **Aggregates**: Consistency boundaries for transactions
6. **Single Responsibility**: Each component has one reason to change
7. **Interface Segregation**: Small, focused interfaces
8. **Dependency Inversion**: Depend on abstractions, not concretions

---

## 🎓 Reference Documentation

- [FILE_STRUCTURE.md](FILE_STRUCTURE.md) - Complete project structure
- [CODE_REVIEW_GUIDELINES.md](CODE_REVIEW_GUIDELINES.md) - Review checklist
- [README.md](../README.md) - Setup and running instructions

---

**When in doubt**: Follow the existing patterns in the codebase. Look at `internal/finance/domain/assetsource/` as a reference implementation.
