# Adapter Layer

This layer contains adapters that convert external requests (HTTP, gRPC, CLI, etc.) into internal domain operations and convert domain responses back to external formats.

## Purpose

- **Protocol Translation**: Convert protocol-specific requests (HTTP, gRPC, etc.) into domain operations
- **Response Formatting**: Format domain responses into the appropriate protocol format (JSON, XML, etc.)
- **Error Handling**: Translate domain errors into protocol-specific error responses
- **Request Validation**: Validate incoming requests before passing to services

## Structure

Each adapter typically handles one protocol or interface type. For example:
- `http/` - HTTP REST endpoints
- `grpc/` - gRPC service definitions
- `cli/` - Command-line interface handlers
- `kafka/` - Kafka message

## Key Principles

- Adapters should be thin and focused on translation
- Business logic should reside in the domain and service layers
- Adapters should not contain validation logic beyond format validation
