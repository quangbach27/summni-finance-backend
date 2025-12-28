# Code Review Guidelines

Senior Golang Developer Code Review Checklist for Sumni Finance Backend

## 🏗️ Architecture & Design

### Clean Architecture Compliance

- [ ] Proper separation of concerns: domain, application, infrastructure layers
- [ ] Domain layer has no external dependencies (no database, HTTP, etc.)
- [ ] Dependencies flow inward (outer layers depend on inner layers)
- [ ] Ports and adapters properly implemented

### SOLID Principles

- [ ] **Single Responsibility**: Each function/struct has one clear purpose
- [ ] **Open/Closed**: Code is open for extension, closed for modification
- [ ] **Liskov Substitution**: Interfaces can be substituted without breaking behavior
- [ ] **Interface Segregation**: Interfaces are small and focused
- [ ] **Dependency Inversion**: Depend on abstractions, not concretions

### Domain-Driven Design

- [ ] Entities properly represent identity and lifecycle
- [ ] Value objects are immutable and validated
- [ ] Aggregates enforce transaction boundaries
- [ ] Repository interfaces in domain layer, implementations in adapter layer
- [ ] Domain logic stays in domain layer (not leaking to handlers/adapters)

---

## 📝 Go Idioms & Best Practices

### Effective Go Compliance

- [ ] Follow standard library patterns
- [ ] Naming conventions: camelCase for unexported, PascalCase for exported
- [ ] Package names: short, concise, lowercase, no underscores
- [ ] Avoid stuttering: `user.UserID` → `user.ID`

### Interface Design

- [ ] Accept interfaces, return structs
- [ ] Interfaces are small and focused (1-3 methods ideal)
- [ ] Interfaces defined by consumer, not producer
- [ ] `-er` naming convention where appropriate (Reader, Writer, Handler)

### Go Fundamentals

- [ ] Proper nil handling (check before dereferencing)
- [ ] Utilize zero values effectively
- [ ] Defer used for resource cleanup
- [ ] Context passed as first parameter
- [ ] Avoid naked returns in long functions

---

## ⚠️ Error Handling

### Error Practices

- [ ] Errors are wrapped with context: `fmt.Errorf("operation failed: %w", err)`
- [ ] Error messages start with lowercase (unless proper noun)
- [ ] Errors provide enough context for debugging
- [ ] Custom error types for domain-specific errors
- [ ] No ignored errors (`_ = err`) without justification
- [ ] Early returns to reduce nesting
- [ ] No panic in library code (only in truly exceptional cases)

### Error Patterns

```go
// ✅ Good
if err != nil {
    return fmt.Errorf("failed to create asset source: %w", err)
}

// ❌ Bad
if err != nil {
    return err  // No context
}

// ❌ Bad
_ = someOperation()  // Ignored error
```

---

## 🔄 Concurrency

### Goroutine Management

- [ ] All goroutines have proper termination paths
- [ ] No goroutine leaks (use context cancellation)
- [ ] WaitGroups used correctly for synchronization
- [ ] Proper channel closing (only sender closes)
- [ ] Channel direction specified where possible

### Race Condition Prevention

- [ ] Tests run with `-race` flag
- [ ] Shared state protected with mutexes or channels
- [ ] Avoid data races in closures (capture variables correctly)
- [ ] Context cancellation respected in all goroutines

### Concurrency Patterns

```go
// ✅ Good: Context cancellation
func process(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case item := <-workCh:
            // process item
        }
    }
}

// ❌ Bad: Goroutine leak
go func() {
    for item := range ch {  // What if ch never closes?
        process(item)
    }
}()
```

---

## 🚀 Performance

### Memory Optimization

- [ ] Minimize unnecessary allocations
- [ ] Use pointers for large structs (>64 bytes)
- [ ] Preallocate slices when size is known: `make([]T, 0, expectedSize)`
- [ ] Preallocate maps when size is known: `make(map[K]V, expectedSize)`
- [ ] Use `strings.Builder` for string concatenation in loops

### Performance Patterns

```go
// ✅ Good: Preallocated slice
items := make([]Item, 0, len(sources))
for _, src := range sources {
    items = append(items, src.ToItem())
}

// ❌ Bad: Repeated allocations
var items []Item
for _, src := range sources {
    items = append(items, src.ToItem())
}

// ✅ Good: strings.Builder
var sb strings.Builder
for _, s := range strings {
    sb.WriteString(s)
}
return sb.String()

// ❌ Bad: Repeated string concatenation
result := ""
for _, s := range strings {
    result += s  // New allocation each iteration
}
```

### Avoid Common Pitfalls

- [ ] No defer in tight loops (performance impact)
- [ ] Avoid unnecessary conversions
- [ ] Use appropriate data structures (map vs slice)

---

## 🧪 Testing

### Test Coverage

- [ ] Unit tests for all business logic
- [ ] Critical paths have test coverage
- [ ] Edge cases tested (nil, empty, boundary values)
- [ ] Error paths tested
- [ ] Integration tests for database operations

### Test Quality

- [ ] Table-driven tests for multiple scenarios
- [ ] Clear test names: `TestFunctionName_Scenario_ExpectedBehavior`
- [ ] Tests are independent (no shared state)
- [ ] Proper setup and teardown
- [ ] Mocks used appropriately (via interfaces)

### Test Patterns

```go
// ✅ Good: Table-driven test
func TestCreateAssetSource(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateInput
        want    *AssetSource
        wantErr bool
    }{
        {
            name: "valid input creates asset source",
            input: CreateInput{Name: "Test", Type: "Bank"},
            want: &AssetSource{Name: "Test", Type: "Bank"},
        },
        {
            name: "empty name returns error",
            input: CreateInput{Name: "", Type: "Bank"},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := CreateAssetSource(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("CreateAssetSource() error = %v, wantErr %v", err, tt.wantErr)
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("CreateAssetSource() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

---

## 🔒 Security

### Input Validation

- [ ] All external inputs validated
- [ ] SQL injection prevented (using SQLC parameterized queries)
- [ ] XSS prevention for user-generated content
- [ ] File upload validation (type, size)

### Authentication & Authorization

- [ ] Proper authentication middleware applied
- [ ] JWT tokens validated correctly
- [ ] Authorization checks for protected resources
- [ ] Session management secure

### Secrets & Configuration

- [ ] No hardcoded secrets, API keys, passwords
- [ ] Environment variables for configuration
- [ ] Sensitive data not logged
- [ ] TLS/HTTPS for all external communication

### Dependency Security

- [ ] Regular dependency updates
- [ ] Scan for vulnerabilities: `go list -m -json all | nancy sleuth`
- [ ] Review third-party dependencies before adding

---

## 📚 Code Organization

### Package Structure

- [ ] Logical package grouping
- [ ] No circular dependencies
- [ ] Package names reflect purpose
- [ ] Avoid generic names (utils, helpers, common for new packages)

### Naming Conventions

- [ ] Clear, descriptive names
- [ ] Avoid abbreviations (unless widely known)
- [ ] Consistent terminology across codebase
- [ ] Boolean names: `isActive`, `hasPermission`, `canDelete`

### Code Structure

- [ ] Functions are focused (20-30 lines ideal)
- [ ] Files are manageable size (<500 lines)
- [ ] Minimize exported surface area
- [ ] Related code grouped together

---

## 📖 Documentation

### Required Documentation

- [ ] Package comment for every package
- [ ] Godoc comments for all exported types, functions, constants
- [ ] Complex logic has inline comments explaining "why", not "what"
- [ ] README updated with new features/changes
- [ ] Architecture decisions documented

### Documentation Standards

```go
// ✅ Good: Comprehensive godoc
// CreateAssetSource creates a new asset source with the given details.
// It validates the input and returns an error if validation fails.
// The created asset source is persisted to the database.
func CreateAssetSource(ctx context.Context, input CreateInput) (*AssetSource, error)

// ❌ Bad: Obvious/redundant comment
// CreateAssetSource creates asset source
func CreateAssetSource(ctx context.Context, input CreateInput) (*AssetSource, error)
```

---

## 🗄️ Database & Persistence

### SQLC Usage

- [ ] All queries defined in `.sql` files
- [ ] Type-safe SQLC generated code used
- [ ] No raw SQL strings in code
- [ ] Proper SQL query optimization

### Transaction Management

- [ ] Transactions used for multi-step operations
- [ ] Proper commit/rollback handling
- [ ] Transaction boundaries aligned with aggregates
- [ ] No long-running transactions

### Query Efficiency

- [ ] No N+1 query problems
- [ ] Appropriate use of joins vs multiple queries
- [ ] Indexes on frequently queried columns
- [ ] LIMIT and pagination for large result sets

### Migrations

- [ ] All schema changes have up/down migrations
- [ ] Migrations are idempotent
- [ ] Migration naming follows convention: `YYYYMMDDHHMMSS_description.up.sql`
- [ ] Backward compatible when possible

---

## 🔧 Dependencies

### Dependency Management

- [ ] Only necessary dependencies added
- [ ] `go.mod` and `go.sum` updated correctly
- [ ] No `+incompatible` versions (upgrade or vendor)
- [ ] Dependencies reviewed before adding

### Standard Library First

- [ ] Prefer standard library over third-party when possible
- [ ] Evaluate necessity of new dependencies
- [ ] Consider maintenance and community support

---

## 🎯 Project-Specific (Sumni Finance)

### CQRS Pattern

- [ ] Commands and queries clearly separated
- [ ] Commands in `app/command/`, queries in `app/query/`
- [ ] Command handlers modify state
- [ ] Query handlers only read state

### Repository Pattern

- [ ] Repository interfaces in domain layer
- [ ] Repository implementations in adapter layer
- [ ] No database-specific code in domain
- [ ] Repository methods align with aggregate operations

### Value Objects

- [ ] Immutable value objects (no setters)
- [ ] Validation in constructor/factory
- [ ] Equality based on value, not identity
- [ ] Examples: Money, Currency, SourceDetails

### Hexagonal Architecture

- [ ] Domain layer is pure (no external dependencies)
- [ ] Ports define interfaces (HTTP handlers)
- [ ] Adapters implement infrastructure (database, external APIs)
- [ ] Application layer orchestrates use cases

---

## 🚨 Common Anti-Patterns to Flag

### Critical Issues

- ❌ Ignoring errors: `_ = operation()`
- ❌ Goroutine leaks (no termination)
- ❌ Race conditions
- ❌ SQL injection vulnerabilities
- ❌ Hardcoded secrets
- ❌ Global mutable state

### Code Smells

- ❌ Deep nesting (arrow code)
- ❌ Long functions (>50 lines)
- ❌ God objects (too many responsibilities)
- ❌ Premature optimization
- ❌ Over-engineering (unnecessary abstractions)
- ❌ Interface pollution (interfaces everywhere)

### Go-Specific

- ❌ Not running `go fmt` / `goimports`
- ❌ Using `init()` for complex initialization
- ❌ Returning interfaces instead of structs
- ❌ Context not as first parameter
- ❌ Mixing tabs and spaces (use tabs)

---

## 🛠️ Automated Checks (CI/CD)

### Required Checks

```bash
# Formatting
go fmt ./...

# Vet (common issues)
go vet ./...

# Linting
golangci-lint run

# Tests with race detector
go test -race -cover ./...

# Benchmarks
go test -bench=. -benchmem

# Security scanning
gosec ./...

# Dependency vulnerabilities
go list -m -json all | nancy sleuth
```

### CI/CD Pipeline

- [ ] All checks pass in CI
- [ ] Tests pass with race detector
- [ ] Linter warnings addressed
- [ ] Code coverage maintained/improved
- [ ] Build succeeds

---

## 📋 Review Process Checklist

### Before Submitting PR

- [ ] Code compiles without errors
- [ ] All tests pass locally
- [ ] `go fmt` applied
- [ ] `golangci-lint` passes
- [ ] Race detector passes
- [ ] Manual testing completed
- [ ] Documentation updated

### During Code Review

1. **First Pass**: Architecture and design

   - Does it fit the project structure?
   - Are layers properly separated?
   - Does it follow CQRS/DDD patterns?

2. **Second Pass**: Code quality

   - Error handling correct?
   - Tests adequate?
   - Performance considerations?
   - Security issues?

3. **Third Pass**: Details
   - Naming clear?
   - Documentation sufficient?
   - Edge cases handled?
   - No anti-patterns?

### Review Etiquette

- 🎯 Focus on the code, not the person
- 💬 Ask questions: "Could we...?" instead of "You should..."
- 📚 Explain rationale, help others learn
- ⚖️ Balance perfection with pragmatism
- ✅ Acknowledge good practices
- 🔴 Differentiate blocking vs non-blocking comments

---

## 💡 Reviewer Guidelines

### Prioritization

1. **Critical**: Security, data corruption, major bugs
2. **Important**: Architecture violations, performance issues, error handling
3. **Nice-to-have**: Style improvements, better naming, minor refactoring

### Comment Format

```markdown
[BLOCKING] SQL injection vulnerability - use parameterized query
[IMPORTANT] Goroutine leak - add context cancellation
[SUGGESTION] Consider using strings.Builder for performance
[NIT] Rename variable `d` to `duration` for clarity
[QUESTION] Why are we using a mutex here instead of a channel?
[PRAISE] Great use of table-driven tests!
```

### When to Approve

- ✅ All blocking issues resolved
- ✅ Tests pass and provide good coverage
- ✅ Architecture aligns with project patterns
- ✅ No security vulnerabilities
- ✅ Documentation adequate

### When to Request Changes

- 🔴 Security vulnerabilities present
- 🔴 Tests missing or failing
- 🔴 Architecture violations
- 🔴 Major bugs or logic errors
- 🔴 No error handling

---

## 🎓 Learning Resources

### Go Best Practices

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

### Architecture

- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Domain-Driven Design by Eric Evans](https://www.domainlanguage.com/ddd/)

### Security

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Go Security Best Practices](https://github.com/Checkmarx/Go-SCP)

---

## 📊 Metrics to Track

### Code Quality Metrics

- Test coverage: Target >80% for critical paths
- Cyclomatic complexity: Keep functions <10
- Code churn: High churn may indicate design issues
- Bug density: Bugs per 1000 lines of code

### Review Metrics

- Time to first review: <24 hours
- Time to merge: <48 hours for small PRs
- Review iteration count: Minimize back-and-forth
- Review thoroughness: Balance speed with quality

---

**Remember**: The goal is to maintain high code quality while fostering a collaborative learning environment. Reviews should be thorough but respectful, focusing on building better software together.
