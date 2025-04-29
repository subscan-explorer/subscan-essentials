.PHONY: check run build

GOCMD=go
BUILD_PATH=cmd
CURRENT_DIR=$(shell pwd)

export CGO_ENABLED=0

getdeps:
	mkdir -p $(GOPATH)/bin
	which golangci-lint 1>/dev/null || (echo "Installing golangci-lint" && $(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5)

lint: getdeps
	echo "Running $@ check"
	${GOPATH}/bin/golangci-lint cache clean
	${GOPATH}/bin/golangci-lint run --timeout=5m --config ./.golangci.yml

check: lint

run:
	$(GOCMD) run ./$(BUILD_PATH)

build:
	$(GOCMD) build -o ./bin/subscan -v ./$(BUILD_PATH)

doc:
	swag init -g cmd/main.go -o ./docs/api --parseInternal --parseDependency
