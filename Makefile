.DEFAULT_GOAL := help

BIN := $(CURDIR)/bin

MOCKGEN_VERSION := v0.6.0

SWAGGER_VERSION := v1.16.2

BUF_VERSION := v1.50.0

PROTOC_GEN_GO_VERSION := v1.36.6

PROTOC_GEN_GO_GRPC_VERSION := v1.5.1

bin/mockgen:
	@echo "Installing mockgen $(MOCKGEN_VERSION) to $(BIN)/mockgen"
	@GOBIN="$(BIN)" go install go.uber.org/mock/mockgen@$(MOCKGEN_VERSION)

.PHONY: mockgen
mockgen: bin/mockgen ## Generate mocks for interfaces
	@PATH="$(BIN):$$PATH" go generate -v ./...

bin/swag:
	@echo "Installing swag $(SWAGGER_VERSION) to $(BIN)/swag"
	@GOBIN="$(BIN)" go install github.com/swaggo/swag/cmd/swag@$(SWAGGER_VERSION)

.PHONY: swag
swag: bin/swag ## Generate Swagger docs
	@PATH="$(BIN):$$PATH" swag init -d internal/adapter/web -g router/router.go -o internal/adapter/web/docs

bin/buf:
	@echo "Installing buf $(BUF_VERSION) to $(BIN)/buf"
	@GOBIN="$(BIN)" go install github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION)

bin/protoc-gen-go:
	@echo "Installing protoc-gen-go $(PROTOC_GEN_GO_VERSION) to $(BIN)/protoc-gen-go"
	@GOBIN="$(BIN)" go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)

bin/protoc-gen-go-grpc:
	@echo "Installing protoc-gen-go-grpc $(PROTOC_GEN_GO_GRPC_VERSION) to $(BIN)/protoc-gen-go-grpc"
	@GOBIN="$(BIN)" go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)

.PHONY: protoc
protoc: bin/buf bin/protoc-gen-go bin/protoc-gen-go-grpc ## Generate gRPC code from proto files
	@PATH="$(BIN):$$PATH" $(BIN)/buf generate

.PHONY: build
build: ## Build all binaries to bin/
	@go build -o $(BIN)/cli ./cmd/cli
	@go build -o $(BIN)/server ./cmd/server

.PHONY: test
test: ## Run tests
	@go test -cover -race ./...

COMPOSE := docker compose -f env/docker-compose.yml

.PHONY: server-up
server-up: ## Start server with postgres in docker
	@$(COMPOSE) up -d --build
	@echo ""
	@echo "HTTP server: http://localhost:8082/app/v1/notes"
	@echo "Swagger UI:  http://localhost:8082/app/v1/swagger/index.html"
	@echo "gRPC server: localhost:50051"

.PHONY: server-down
server-down: ## Stop server and postgres
	@$(COMPOSE) down

.PHONY: server-reset
server-reset: ## Stop server and postgres, remove volumes
	@$(COMPOSE) down -v

.PHONY: help
help: ## Show available commands
	@grep -E '^\S+:.*##' $(MAKEFILE_LIST) | sed 's/:.* ## / — /' | sort