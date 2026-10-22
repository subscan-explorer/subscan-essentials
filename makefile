.PHONY: check run build

GOCMD=go
GOLANGCI_LINT_VERSION=2.12.2
GOLANGCI_LINT?=$(GOPATH)/bin/golangci-lint
BUILD_PATH=cmd
CURRENT_DIR=$(shell pwd)

# substrate-api-rpc/keyring imports C, so lint and builds require cgo.
export CGO_ENABLED=1

getdeps:
	@$(GOLANGCI_LINT) --version 2>/dev/null | grep -q "version $(GOLANGCI_LINT_VERSION) " || ( \
		mkdir -p "$(dir $(GOLANGCI_LINT))" && \
		GOBIN="$(abspath $(dir $(GOLANGCI_LINT)))" $(GOCMD) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v$(GOLANGCI_LINT_VERSION) \
	)

lint: getdeps
	echo "Running $@ check"
	$(GOLANGCI_LINT) run --timeout=5m --config ./.golangci.yml

check: lint

run:
	$(GOCMD) run ./$(BUILD_PATH)

build:
	$(GOCMD) build -o ./bin/subscan -v ./$(BUILD_PATH)

doc:
	swag init -g cmd/main.go -o ./docs/api --parseInternal --parseDependency
