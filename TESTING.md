# Testing Guide

This document describes the comprehensive testing strategy for the Go Client Locations microservice application.

## Test Structure

The testing is organized into three main categories:

1. **Unit Tests** - Test individual components in isolation
2. **Functional Tests** - Test API endpoints and business logic
3. **Integration Tests** - Test the complete system with real dependencies

## Test Categories

### Unit Tests

Unit tests focus on testing individual functions and methods in isolation using mocks.

**Coverage:**
- Shared utilities (`shared/util/`)
- User service business logic (`services/user-service/internal/service/`)
- Location history service business logic (`services/location-history-service/`)

**Key Features:**
- Table-driven test patterns
- Comprehensive edge case coverage
- Mock dependencies
- Fast execution

### Functional Tests

Functional tests verify API endpoints and business logic with mocked external dependencies.

**Coverage:**
- API Gateway HTTP endpoints
- Request/response validation
- Error handling
- Business logic flows

**Key Features:**
- HTTP request/response testing
- JSON serialization/deserialization
- Input validation
- Error scenarios

### Integration Tests

Integration tests verify the complete system with real dependencies.

**Coverage:**
- End-to-end API flows
- Database operations
- Service communication
- Real system behavior

**Key Features:**
- Full system testing
- Real database connections
- Service-to-service communication
- Performance testing

## Running Tests

### Quick Start

```bash
# Run all unit and functional tests
make test

# Run only unit tests
make test-unit

# Run only functional tests
make test-functional

# Run integration tests (requires full system)
make test-integration

# Run all tests including integration
make test-all
```

### Using the Test Runner Script

```bash
# Run unit tests
./scripts/run-tests.sh -m unit

# Run functional tests
./scripts/run-tests.sh -m functional

# Run integration tests
./scripts/run-tests.sh -m integration

# Run all tests with coverage
./scripts/run-tests.sh -m all -c

# Run with verbose output
./scripts/run-tests.sh -m unit -v
```

### Test Coverage

```bash
# Generate coverage report
make test-coverage

# View coverage in browser
open coverage.html
```

## Test Configuration

### Environment Variables

```bash
# For unit and functional tests
export TEST_MODE=unit
export MOCK_DB=true

# For integration tests
export TEST_MODE=integration
export MONGODB_URI=mongodb://localhost:27017/test
export API_GATEWAY_URL=http://localhost:8004
export USER_SERVICE_URL=localhost:9093
export LOCATION_SERVICE_URL=localhost:9092
```

### Test Data

Test data is managed through the `shared/testutil` package:

```go
// Create test users
testUsers := []*domain.UserModel{
    testutil.CreateTestUser("user1", 51.0, 16.0),
    testutil.CreateTestUser("user2", 52.0, 17.0),
}

// Create test coordinates
coords := testutil.CreateTestCoordinate(51.0, 16.0)
```

## Test Patterns

### Table-Driven Tests

All tests use the table-driven test pattern for comprehensive coverage:

```go
func TestValidateUserName(t *testing.T) {
    tests := []struct {
        name        string
        username    string
        expectError bool
        errorMsg    string
    }{
        {
            name:        "valid username",
            username:    "testuser123",
            expectError: false,
        },
        {
            name:        "invalid username - too short",
            username:    "ab",
            expectError: true,
            errorMsg:    "username must be between 4 and 16 characters",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateUserName(tt.username)
            // Test assertions...
        })
    }
}
```

### Mock Dependencies

Use the `shared/testutil` package for mocking:

```go
// Mock repository
mockRepo := testutil.NewMockUserRepository()
mockRepo.SetUsers(testUsers)
service := NewService(mockRepo)
```

## Test Coverage Requirements

- **Unit Tests**: Minimum 80% coverage
- **Functional Tests**: All endpoints covered
- **Integration Tests**: Critical paths covered

## Continuous Integration

### GitHub Actions

```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.21
      - name: Run unit tests
        run: make test-unit
      - name: Run functional tests
        run: make test-functional
      - name: Run integration tests
        run: make test-integration
        env:
          MONGODB_URI: mongodb://localhost:27017/test
```

## Troubleshooting

### Common Issues

1. **Integration tests failing**: Ensure all services are running
2. **Database connection errors**: Check MongoDB connection
3. **Port conflicts**: Verify service ports are available

### Debug Mode

```bash
# Run tests with debug output
go test -v -timeout=60s ./...

# Run specific test
go test -v -run TestValidateUserName ./shared/util/
```

## Best Practices

1. **Test Isolation**: Each test should be independent
2. **Clear Test Names**: Use descriptive test names
3. **Comprehensive Coverage**: Test happy path and edge cases
4. **Mock External Dependencies**: Use mocks for unit tests
5. **Real Dependencies for Integration**: Use real services for integration tests
6. **Performance Testing**: Include performance tests for critical paths

## Test Data Management

### Test Users

```go
// Predefined test users with known coordinates
testUsers := []*domain.UserModel{
    {
        UserName: "user1",
        Coordinates: &types.Coordinate{
            Latitude:  51.11822470712269,
            Longitude: 16.990711729269563,
        },
    },
    // ... more users
}
```

### Test Coordinates

```go
// Known test coordinates for distance calculations
wroclaw := &types.Coordinate{Latitude: 51.11822470712269, Longitude: 16.990711729269563}
warsaw := &types.Coordinate{Latitude: 52.23553956649786, Longitude: 20.984595191389918}
```

## Performance Testing

### Benchmark Tests

```bash
# Run benchmarks
make bench

# Run specific benchmark
go test -bench=BenchmarkCalculateDistance ./shared/util/
```

### Load Testing

```bash
# Run concurrent tests
go test -race ./...
```

## Test Maintenance

1. **Update Tests**: Keep tests in sync with code changes
2. **Remove Obsolete Tests**: Clean up unused tests
3. **Add New Tests**: Add tests for new features
4. **Review Coverage**: Regularly review test coverage

## Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Table-Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [Mocking in Go](https://github.com/golang/mock)
- [Integration Testing Best Practices](https://martinfowler.com/articles/practical-test-pyramid.html)
