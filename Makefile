BIN := $(CURDIR)/bin

bin/mockgen:
	@mkdir -p bin
	@GOBIN="$(BIN)" go install go.uber.org/mock/mockgen@v0.6.0

.PHONY: mockgen
mockgen: bin/mockgen
	@PATH="$(BIN):$$PATH" go generate -v ./...

.PHONY: test
test:
	@go test -cover -race ./...