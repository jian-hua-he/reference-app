# Application Layer

This layer contains application services (use cases) that orchestrate domain logic and coordinate interactions between the domain layer and external systems. It implements the business use cases of the application.

## Purpose

- **Use Case Implementation**: Each application service implements a specific user workflow or business use case
- **Orchestration**: Coordinate multiple domain operations to fulfill a complete use case

## Code Organization Guidelines

### Service Structure

Each application service should:

1. **Have a single responsibility**: Focus on one aggregate or related set of operations
2. **Depend on interfaces**: Use dependency injection for repositories and other services
3. **Use constructor injection**: Provide a `New<Service>` constructor function
4. **Define dependencies**: Create a `dependency.go` file with all required interfaces
5. **Implement use cases**: Each method represents a complete user workflow

### File Organization

For each service (e.g., NoteApp):
- `note.go` - Service struct and methods
- `dependency.go` - Interface definitions (must match what external layers provide)
- `note_test.go` - Comprehensive test cases
- `dependency_mock.go` - Auto-generated mocks (via mockgen)

### Naming Conventions

- Service types: `<Domain>App` (e.g., `NoteApp`, `UserApp`)
- Constructor: `New<Domain>App` (e.g., `NewNoteApp`)
- Methods: Verb-based names describing the action (e.g., `Create()`, `List()`, `Delete()`)
- Interfaces: `<Domain>Repository`, `<Domain>Service` (e.g., `NoteRepository`)

## Key Principles

- **Thin Controllers**: Adapters should delegate business logic to application services
- **Interface Segregation**: Define only the dependencies you need
- **Error Wrapping**: Wrap lower-level errors with context using `fmt.Errorf`
- **No Framework Coupling**: Application services should not depend on HTTP, database, or other frameworks
- **Testability**: Should be fully testable using mocks for dependencies
- **Immutability**: Return new entities rather than modifying input parameters
