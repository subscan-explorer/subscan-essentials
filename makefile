GOCMD=go
BUILD_PATH=cmd

export GO111MODULE=on

build:
	$(GOCMD) build -o ./cmd/subscan -v ./$(BUILD_PATH)