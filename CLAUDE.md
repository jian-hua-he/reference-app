# Development Workflow

## Implementation Cycle (TDD)

Follow this cycle for every new feature or change:

1. **Create empty struct or function** — Define the signature with a minimal or zero-value implementation.
2. **Write tests and fail them** — Write test cases that exercise the new code. Run `make test` and confirm they fail.
3. **Update implementation to pass the tests** — Fill in the logic until all tests pass.
4. **Repeat** until all requirements are complete.

## Commands

- `make test` — Run all tests with coverage and race detection (`go test -cover -race ./...`).
- `make mockgen` — Regenerate all mocks via `go generate ./...`.
- `make swag` — Regenerate Swagger docs.
- `make protoc` — Regenerate gRPC code from proto files.

## Testing Conventions

### Mocking with gomock

- Use `go.uber.org/mock/mockgen` to generate mocks.
- Each package declares its dependencies as interfaces in a file named `dependencies.go`.
- The `//go:generate` directive in `dependencies.go` produces `dependencies_mock_test.go` with `_test` package suffix.
- Run `make mockgen` to regenerate all mocks after changing interfaces.

### Interfaces on the consumer side

- Define interfaces where they are consumed, not where they are implemented.
- Place them in `dependencies.go` within the consuming package.

### Mock file naming

- Mock file: `dependencies_mock_test.go`
- The mock file package name must always use the `_test` suffix (e.g., `package note_test`).

### Test package

- Test files must always use the `_test` package suffix (e.g., `package note_test`, not `package note`).

### Assertions with testify

- Use `github.com/stretchr/testify/assert` for assertions (`assert.Equal`, `assert.ErrorIs`, etc.).
- Use `require` variants when the test cannot continue after a failure.

### Table-driven tests

- Use `map[string]struct{...}` for table tests so cases are keyed by name:

```go
testCases := map[string]struct {
    Input   string
    Want    string
    WantErr error
}{
    "case name": {
        Input: "value",
        Want:  "expected",
    },
}

for name, tc := range testCases {
    t.Run(name, func(t *testing.T) {
        // ...
    })
}
```

### Integration tests with testcontainers

- When a test requires an external service (database, cache, message broker, etc.), use `testcontainers-go` to spin up the dependency in a container.
- Do not require developers to run external services manually for tests.

## Project Structure

This project follows Clean Architecture with these layers:

- `internal/entity/` — Domain models (pure data structs).
- `internal/repository/` — Data persistence interfaces and implementations.
- `internal/application/` — Use-case services (business logic).
- `internal/adapter/` — Protocol adapters (HTTP, gRPC, CLI).
- `pkg/` — Reusable packages.
- `cmd/` — Application entry points.
