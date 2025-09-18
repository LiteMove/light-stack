# Makefile for light-stack

.PHONY: help build run dev migrate clean test

# Default target
help:
	@echo "Available targets:"
	@echo "  build    - Build the application"
	@echo "  run      - Run the application"
	@echo "  dev      - Run in development mode"
	@echo "  migrate  - Run database migration"
	@echo "  clean    - Clean build files"
	@echo "  test     - Run tests"
	@echo "  web-dev  - Start frontend development server"
	@echo "  web-build- Build frontend for production"

# Build the application
build:
	go mod tidy
	go build -o bin/server cmd/server/main.go
	go build -o bin/migrate cmd/migrate/main.go

# Run the application
run: build
	./bin/server

# Run in development mode with auto-reload
dev:
	go run cmd/server/main.go

# Run database migration
migrate:
	go run cmd/migrate/main.go

# Clean build files
clean:
	rm -rf bin/
	rm -rf web/dist/

# Run tests
test:
	go test ./...

# Frontend development
web-dev:
	cd web && npm run dev

# Build frontend
web-build:
	cd web && npm run build

# Install dependencies
deps:
	go mod download
	cd web && npm install

# Format code
fmt:
	go fmt ./...
	cd web && npm run format

# Lint code
lint:
	golangci-lint run
	cd web && npm run lint