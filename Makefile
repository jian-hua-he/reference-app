BIN := $(CURDIR)/bin

bin/mockgen:
	@mkdir -p bin
	@GOBIN="$(BIN)" go install go.uber.org/mock/mockgen@v0.6.0

mockgen: bin/mockgen
	@PATH="$(BIN):$$PATH" go generate -v ./...