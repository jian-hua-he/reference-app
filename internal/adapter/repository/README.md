# Repository Layer

This layer provides abstractions for data persistence and retrieval. It defines interfaces that represent ways to store and retrieve domain entities.

## Purpose

- **Data Abstraction**: Define interfaces for saving and retrieving domain entities
- **Multiple Implementations**: Support different data storage mechanisms (in-memory, database, file system, etc.)
- **Repository Pattern**: Provide repository interfaces that services depend on

## Structure

- `/${repo_name}/memory/`
- `/${repo_name}/postgres/`
- `/${repo_name}/file/`
- `/${repo_name}/mongodb`

## Key Principles

- Repository interfaces should be defined based on domain concepts, not database operations
- Implementations are interchangeable (swap in-memory for database without changing service code)
- Repositories should only handle persistence concerns, not business logic
