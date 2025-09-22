# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
PROTO_DIR := proto
PROTO_SRC := $(wildcard $(PROTO_DIR)/*.proto)
GO_OUT := .

# Build parameters
BINARY_NAME=api-gateway
BINARY_UNIX=$(BINARY_NAME)_unix

# Test parameters
TEST_TIMEOUT=30s
COVERAGE_OUT=coverage.out
COVERAGE_HTML=coverage.html

.PHONY: all build clean test test-unit test-functional test-integration test-coverage deps

all: deps build

deps:
	$(GOMOD) download
	$(GOMOD) tidy

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./services/api-gateway
	$(GOBUILD) -o location-history-service -v ./services/location-history-service
	$(GOBUILD) -o user-service -v ./services/user-service/cmd

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f location-history-service
	rm -f user-service
	rm -f $(COVERAGE_OUT)
	rm -f $(COVERAGE_HTML)

# Run all tests (unit and functional only)
test: test-unit test-functional

# Run only unit tests
test-unit:
	@echo "Running unit tests..."
	$(GOTEST) -v -timeout=$(TEST_TIMEOUT) ./shared/util/...
	$(GOTEST) -v -timeout=$(TEST_TIMEOUT) ./services/user-service/internal/service/...
	$(GOTEST) -v -timeout=$(TEST_TIMEOUT) ./services/location-history-service/...

# Run functional tests (API Gateway endpoints)
test-functional:
	@echo "Running functional tests..."
	@echo "Note: API Gateway functional tests require gRPC services to be running"
	@echo "These tests are included in integration tests instead"
	@echo "Functional tests completed (API Gateway tests skipped - require full system)"

# Run integration tests (requires full system)
test-integration:
	@echo "Running integration tests..."
	$(GOTEST) -v -timeout=60s -tags=integration ./services/api-gateway/...
	$(GOTEST) -v -timeout=60s -tags=integration ./services/user-service/internal/infrastructure/repository/...

# Run all tests including integration
test-all: test test-integration

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -timeout=$(TEST_TIMEOUT) -coverprofile=$(COVERAGE_OUT) ./...
	$(GOTEST) -v -timeout=$(TEST_TIMEOUT) -coverprofile=$(COVERAGE_OUT) -tags=integration ./...
	$(GOCMD) tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	$(GOTEST) -v -timeout=$(TEST_TIMEOUT) -race ./...

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -v -bench=. -benchmem ./...

# Run specific test patterns
test-pattern:
	@echo "Running tests matching pattern: $(PATTERN)"
	$(GOTEST) -v -timeout=$(TEST_TIMEOUT) -run=$(PATTERN) ./...

# Lint code
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	$(GOCMD) vet ./...

# Security scan
security:
	@echo "Running security scan..."
	gosec ./...

# Generate test coverage badge
coverage-badge:
	@echo "Generating coverage badge..."
	$(GOTEST) -v -coverprofile=$(COVERAGE_OUT) ./...
	$(GOCMD) tool cover -func=$(COVERAGE_OUT) | tail -1 | awk '{print $$3}' | sed 's/%//' > coverage.txt

# Generate proto
generate-proto:
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(GO_OUT) \
		--go-grpc_out=$(GO_OUT) \
		$(PROTO_SRC)

# Help target
help:
	@echo "Available targets:"
	@echo "  all              - Build all services"
	@echo "  build            - Build all services"
	@echo "  clean            - Clean build artifacts"
	@echo "  test             - Run unit and functional tests"
	@echo "  test-unit        - Run only unit tests"
	@echo "  test-functional  - Run functional tests"
	@echo "  test-integration - Run integration tests (requires full system)"
	@echo "  test-all         - Run all tests including integration"
	@echo "  test-coverage    - Run tests with coverage report"
	@echo "  test-race        - Run tests with race detection"
	@echo "  bench            - Run benchmarks"
	@echo "  test-pattern     - Run tests matching pattern (use PATTERN=pattern)"
	@echo "  lint             - Run linter"
	@echo "  fmt              - Format code"
	@echo "  vet              - Vet code"
	@echo "  security         - Run security scan"
	@echo "  coverage-badge    - Generate coverage badge"
	@echo "  help             - Show this help message"