# LOCAL_BIN_DIR := $(PWD)/bin
# MOCKGEN_VERSION := v0.5.0

# bin/mockgen:
# 	@echo "> downloading mockgen@$(MOCKGEN_VERSION) to $(LOCAL_BIN_DIR)"
# 	@GOBIN="$(LOCAL_BIN_DIR)" go install go.uber.org/mock/mockgen@$(MOCKGEN_VERSION)

# .PHONY: mockgen
# mockgen: bin/mockgen
# 	@PATH="$(LOCAL_BIN_DIR):$(PATH)" go generate -v ./...

BIN := $(CURDIR)/bin

bin/mockgen:
	@mkdir -p bin
	@GOBIN="$(BIN)" go install go.uber.org/mock/mockgen@v0.6.0

mockgen: bin/mockgen
	@PATH="$(BIN):$$PATH" go generate ./...