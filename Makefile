BINARY_NAME := gpp

.PHONY: build run fmt test lint

build:
	go build -o bin/$(BINARY_NAME) ./cmd/gpp

run:
	go run ./cmd/gpp

fmt:
	go fmt ./...

test:
	go test ./...

lint:
	@echo "lint placeholder (add golangci-lint if needed)"
