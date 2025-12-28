# Sumni Finance Backend - File Structure

This document provides a comprehensive overview of the project's file structure, organized by domain and responsibility.

## 🏗️ Overall Architecture

The project follows **Clean Architecture** principles with **Domain-Driven Design (DDD)** patterns, organized into three main domains:

- **Finance Domain**: Core business logic for financial operations
- **Common Domain**: Shared utilities and cross-cutting concerns
- **Auth Domain**: Authentication and authorization (Keycloak integration)

---

## 📁 Root Level

```
sumni-finance-backend/
├── cmd/                          # Application entry points
├── db/                           # Database migrations
├── docker/                       # Docker configuration files
├── internal/                     # Internal application code
├── realms/                       # Keycloak realm configurations
├── scripts/                      # Utility scripts
├── .github/                      # GitHub workflows and configurations
├── docker-compose.yml            # Docker Compose orchestration
├── go.mod                        # Go module dependencies
├── go.sum                        # Go module checksums
├── Makefile                      # Build and development tasks
├── lefthook.yml                  # Git hooks configuration
├── .golangci.yml                 # Go linter configuration
├── .mockery.yml                  # Mock generation configuration
└── README.md                     # Project documentation
```

---

## 🚀 Entry Points (`cmd/`)

```
cmd/
└── server/
    └── main.go                   # Main application entry point
                                  # - Initializes logging
                                  # - Sets up database connection pool
                                  # - Configures HTTP routes
                                  # - Integrates all domains
```

**Purpose**: Application bootstrapping and dependency injection.

---

## 🗄️ Database (`db/`)

```
db/
└── migrations/
    ├── 000001_init_asset_source.up.sql      # Asset source table creation
    └── 000001_init_asset_source.down.sql    # Asset source table rollback
```

**Purpose**: Database schema migrations managed via migration tool (likely golang-migrate).

---

## 🐳 Docker Configuration (`docker/`)

```
docker/
├── app-local/
│   ├── Dockerfile                # Development Dockerfile
│   ├── reflex.conf              # Hot-reload configuration
│   └── start.sh                 # Container startup script
└── keycloak/
    └── SumniFinanceApp-realm.json    # Keycloak realm export
```

**Purpose**:

- Local development environment with hot-reload
- Keycloak configuration for authentication
- Optional debugging setup with Delve

---

## 💼 Finance Domain (`internal/finance/`)

The **Finance Domain** follows Hexagonal Architecture (Ports & Adapters) with CQRS pattern.

```
internal/finance/
├── domain/                       # Core business logic (Enterprise Business Rules)
│   ├── assetsource/
│   │   ├── asset_source.go      # Asset source entity
│   │   ├── asset_source_test.go
│   │   ├── source_details.go    # Value object for source details
│   │   ├── source_details_test.go
│   │   └── repository.go        # Repository interface
│   ├── transaction/             # Transaction aggregate (future)
│   └── wallet/
│       ├── wallet.go            # Wallet entity
│       ├── allocation.go        # Allocation value object
│       └── repository.go        # Wallet repository interface
│
├── app/                         # Application layer (Application Business Rules)
│   ├── app.go                   # Application struct defining Commands & Queries
│   ├── command/
│   │   └── create_asset_source.go    # Command handler for creating asset sources
│   └── query/
│       ├── get_asset_source.go       # Query handler for retrieving asset sources
│       └── types.go                  # Query DTOs
│
├── adapter/                     # Infrastructure adapters (Interface Adapters)
│   ├── db/
│   │   ├── assetsource_repo.go       # Asset source repository implementation
│   │   └── store/                    # SQLC generated code
│   │       ├── sqlc.yml             # SQLC configuration
│   │       ├── sqlc.go              # SQLC database connection
│   │       ├── models.go            # Generated models
│   │       ├── asset_sources.sql.go # Generated queries
│   │       └── queries/
│   │           └── asset_sources.sql # SQL query definitions
│   └── intra_process/               # In-process communication (future)
│
├── ports/                       # HTTP interface (External Interface)
│   ├── http.go                  # HTTP routes registration
│   └── asset_source_handler.go # HTTP handlers for asset source endpoints
│
└── service/
    └── application.go           # Application dependency injection
```

### Finance Domain Layers

1. **Domain Layer** (`domain/`)

   - Pure business logic with no external dependencies
   - Entities, Value Objects, and Repository interfaces
   - Aggregates: AssetSource, Transaction (future), Wallet

2. **Application Layer** (`app/`)

   - CQRS pattern: Commands for writes, Queries for reads
   - Use cases and application business rules
   - Orchestrates domain objects

3. **Adapter Layer** (`adapter/`)

   - Database implementations using SQLC for type-safe queries
   - External service integrations
   - Infrastructure concerns

4. **Ports Layer** (`ports/`)

   - HTTP handlers and API endpoints
   - Request/Response DTOs
   - API contract definition

5. **Service Layer** (`service/`)
   - Application composition and dependency injection
   - Wiring up all layers

---

## 🔧 Common Domain (`internal/common/`)

Shared utilities and cross-cutting concerns used across all domains.

```
internal/common/
├── cqrs/                        # CQRS infrastructure
│   ├── command.go              # Base command interfaces
│   ├── query.go                # Base query interfaces
│   └── logging.go              # CQRS logging middleware
│
├── db/                         # Database utilities
│   ├── pgx_connection.go       # PostgreSQL connection pool setup
│   ├── transaction.go          # Transaction management utilities
│   └── covert_pgtype.go        # pgx type conversions
│
├── logs/                       # Logging infrastructure
│   ├── init.go                 # Logger initialization
│   └── http.go                 # HTTP request logging middleware
│
├── server/                     # HTTP server utilities
│   ├── http.go                 # HTTP server setup and configuration
│   └── httperr/
│       ├── http_err.go         # HTTP error types
│       └── errors.go           # Error handling utilities
│
├── validator/                  # Request validation
│   ├── validator.go            # Validation logic
│   ├── validator_test.go
│   └── errors_list.go          # Validation error formatting
│
└── valueobject/                # Shared value objects
    ├── money.go                # Money value object
    ├── money_test.go
    ├── currency.go             # Currency value object
    └── currency_test.go
```

### Common Domain Purpose

- **CQRS**: Command and Query abstraction for consistent application layer
- **Database**: Connection pooling and transaction management
- **Logging**: Structured logging with HTTP request context
- **HTTP Server**: Server setup, routing, and error handling
- **Validation**: Request validation with error formatting
- **Value Objects**: Shared domain primitives (Money, Currency)

---

## 🔐 Auth Domain (`internal/auth/`)

Authentication and authorization using **Keycloak** as the identity provider.

```
internal/auth/
├── keycloak_client.go                  # Keycloak OIDC client
│                                       # - OAuth2/OIDC integration
│                                       # - Token verification
│                                       # - Authorization code flow with PKCE
│                                       # - Token refresh logic
│
├── auth_http.go                        # HTTP authentication handlers
│                                       # - Login endpoint
│                                       # - Callback endpoint
│                                       # - Logout endpoint
│                                       # - Auth middleware
│
├── token_inmemory_respository.go       # In-memory token storage
│                                       # - Session management
│                                       # - Token caching
│
└── token_inmemory_repository_test.go   # Token repository tests
```

### Auth Domain Features

- **OAuth2/OIDC Integration**: Full OAuth2 Authorization Code Flow with PKCE
- **Token Management**: Token storage, refresh, and verification
- **Keycloak Integration**: Uses Keycloak as authorization server
- **Middleware**: Authentication middleware for protected routes
- **Session Management**: In-memory session storage (can be replaced with Redis)

**Note**: Currently commented out in main.go but ready to be enabled.

---

## ⚙️ Configuration (`internal/config/`)

```
internal/config/
├── config.go                   # Configuration loading from environment
└── config_test.go              # Configuration tests
```

**Purpose**: Centralized configuration management using environment variables.

---

## 🧪 Scripts (`scripts/`)

```
scripts/
└── test.sh                     # Test execution script
```

---

## 🔁 Development Workflow

### Key Commands (from Makefile)

- `make dev` - Start development environment with hot-reload
- `make dev DEBUG=true` - Start with debugging enabled (Delve on port 40000)
- `make stop` - Stop all containers
- `make test` - Run tests

### Hot Reload

- Uses **reflex** for automatic recompilation on file changes
- Configuration in `docker/app-local/reflex.conf`

### Debugging

- **Delve** debugger available on port 40000
- VS Code debugging configuration in `.vscode/launch.json`

---

## 🏛️ Architecture Patterns

### 1. **Clean Architecture**

- Dependency rule: inner layers don't depend on outer layers
- Domain layer is pure business logic
- Infrastructure details at the edges

### 2. **Hexagonal Architecture (Ports & Adapters)**

- **Domain**: Core business logic
- **Ports**: Interfaces to the outside world (HTTP handlers)
- **Adapters**: Concrete implementations (database, external services)

### 3. **Domain-Driven Design (DDD)**

- **Entities**: Asset Source, Wallet
- **Value Objects**: Money, Currency, SourceDetails, Allocation
- **Repositories**: Data access abstractions
- **Aggregates**: Bounded consistency boundaries

### 4. **CQRS (Command Query Responsibility Segregation)**

- **Commands**: Write operations (CreateAssetSource)
- **Queries**: Read operations (GetAssetSource)
- Separate models for reads and writes

### 5. **Dependency Injection**

- Manual DI in `service/application.go`
- Constructor functions for dependency wiring

---

## 🗃️ Data Layer

### SQLC

- Type-safe SQL query generation
- Configuration in `internal/finance/adapter/db/store/sqlc.yml`
- SQL queries in `queries/` directory
- Generated code in `store/` directory

### Database Migrations

- Located in `db/migrations/`
- Versioned migration files
- Up/down migrations for schema changes

### Connection Pool

- pgx v5 connection pooling
- Configuration in `internal/common/db/pgx_connection.go`

---

## 🔒 Security

- **Keycloak**: Enterprise-grade identity and access management
- **OAuth2/OIDC**: Standard authentication protocols
- **PKCE**: Protection against authorization code interception
- **Token Verification**: JWT signature verification
- **Middleware**: Route protection with authentication middleware

---

## 📝 Testing Strategy

- Unit tests for domain logic (entities, value objects)
- Repository tests for data access
- Configuration tests
- Test scripts in `scripts/test.sh`
- Test environment variables in `.e2e.env`

---

## 🛠️ Tooling

- **SQLC**: Type-safe SQL query generation
- **golang-migrate**: Database migrations
- **reflex**: Hot reload for development
- **Delve**: Go debugger
- **lefthook**: Git hooks management
- **golangci-lint**: Go linting
- **mockery**: Mock generation for testing

---

## 🚦 Future Domains/Features

Based on the current structure, planned future implementations:

1. **Transaction Domain** (`internal/finance/domain/transaction/`)

   - Financial transaction management
   - Transaction aggregates

2. **Wallet Domain Expansion** (`internal/finance/domain/wallet/`)

   - Wallet management
   - Asset allocation

3. **Intra-process Communication** (`internal/finance/adapter/intra_process/`)
   - Domain event handling
   - Event-driven architecture

---

## 📚 Key Design Decisions

1. **Monolithic Modular Architecture**: Single deployable with clear domain boundaries
2. **CQRS**: Separation of read and write operations for scalability
3. **SQLC over ORM**: Type-safe SQL with full control over queries
4. **Keycloak**: External authentication for security and SSO capabilities
5. **Docker-first Development**: Consistent development environment
6. **Hot Reload**: Fast feedback loop during development

---

## 🎯 Domain Boundaries

### Finance Domain

- **Responsibility**: Financial operations, asset management, transactions
- **Entry Point**: `internal/finance/ports/http.go`
- **Core Entities**: AssetSource, Wallet, Transaction

### Auth Domain

- **Responsibility**: Authentication, authorization, session management
- **Entry Point**: `internal/auth/auth_http.go`
- **Integration**: Keycloak OIDC provider

### Common Domain

- **Responsibility**: Shared utilities, infrastructure concerns
- **Usage**: Imported by all other domains
- **Scope**: Cross-cutting concerns only

---

This structure supports:

- ✅ Testability
- ✅ Maintainability
- ✅ Scalability
- ✅ Clear separation of concerns
- ✅ Domain isolation
- ✅ Flexibility for future growth
