# Testing Guide


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

