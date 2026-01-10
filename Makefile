BIN := $(CURDIR)/bin

MOCKGEN_VERSION := v0.6.0

SWAGGER_VERSION := v1.16.2

bin/mockgen:
	@echo "Installing mockgen $(MOCKGEN_VERSION) to $(BIN)/mockgen"
	@GOBIN="$(BIN)" go install go.uber.org/mock/mockgen@$(MOCKGEN_VERSION)

.PHONY: mockgen
mockgen: bin/mockgen
	@PATH="$(BIN):$$PATH" go generate -v ./...

bin/swag:
	@echo "Installing swag $(SWAGGER_VERSION) to $(BIN)/swag"
	@GOBIN="$(BIN)" go install github.com/swaggo/swag/cmd/swag@$(SWAGGER_VERSION)

.PHONY: swag
swag: bin/swag
	@PATH="$(BIN):$$PATH" swag init -d internal/adapter/web -g router/router.go -o internal/adapter/web/docs

.PHONY: test
test:
	@go test -cover -race ./...