# Domain Layer

This layer contains the core business logic and domain models that represent the entities and concepts central to the application.

## Purpose

- **Business Entities**: Define domain objects (e.g., `Note`, `User`) that represent business concepts
- **Business Rules**: Encode validation rules and constraints that must always be true
- **Value Objects**: Define immutable objects that represent domain concepts (e.g., `ID`, `Email`)

## Key Principles

- Domain entities should be independent of framework or external dependencies
- Domain logic should be testable without requiring infrastructure
- Entities should enforce their own invariants and business rules
- No dependencies on repositories, services, or adapters
