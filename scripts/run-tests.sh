#!/bin/bash

# Test Runner Script
# This script provides different test execution modes

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
MODE="unit"
VERBOSE=false
COVERAGE=false
RACE=false
INTEGRATION=false

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -m, --mode MODE        Test mode: unit, functional, integration, all (default: unit)"
    echo "  -v, --verbose          Verbose output"
    echo "  -c, --coverage         Generate coverage report"
    echo "  -r, --race             Run with race detection"
    echo "  -i, --integration      Include integration tests"
    echo "  -h, --help             Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                     # Run unit tests"
    echo "  $0 -m functional       # Run functional tests"
    echo "  $0 -m all -c           # Run all tests with coverage"
    echo "  $0 -m integration -i   # Run integration tests"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -m|--mode)
            MODE="$2"
            shift 2
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -c|--coverage)
            COVERAGE=true
            shift
            ;;
        -r|--race)
            RACE=true
            shift
            ;;
        -i|--integration)
            INTEGRATION=true
            shift
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Validate mode
case $MODE in
    unit|functional|integration|all)
        ;;
    *)
        print_error "Invalid mode: $MODE"
        show_usage
        exit 1
        ;;
esac

# Build test command
TEST_CMD="go test"

# Add verbose flag
if [ "$VERBOSE" = true ]; then
    TEST_CMD="$TEST_CMD -v"
fi

# Add race detection
if [ "$RACE" = true ]; then
    TEST_CMD="$TEST_CMD -race"
fi

# Add coverage
if [ "$COVERAGE" = true ]; then
    TEST_CMD="$TEST_CMD -coverprofile=coverage.out"
fi

# Add timeout
TEST_CMD="$TEST_CMD -timeout=30s"

# Function to run unit tests
run_unit_tests() {
    print_status "Running unit tests..."
    
    # Test shared utilities
    $TEST_CMD ./shared/util/...
    if [ $? -eq 0 ]; then
        print_success "Shared utilities tests passed"
    else
        print_error "Shared utilities tests failed"
        return 1
    fi
    
    # Test user service
    $TEST_CMD ./services/user-service/internal/service/...
    if [ $? -eq 0 ]; then
        print_success "User service tests passed"
    else
        print_error "User service tests failed"
        return 1
    fi
    
    # Test location history service
    $TEST_CMD ./services/location-history-service/...
    if [ $? -eq 0 ]; then
        print_success "Location history service tests passed"
    else
        print_error "Location history service tests failed"
        return 1
    fi
}

# Function to run functional tests
run_functional_tests() {
    print_status "Running functional tests..."
    print_warning "API Gateway functional tests require gRPC services to be running"
    print_warning "Skipping API Gateway functional tests for now"
    
    # Note: API Gateway tests require gRPC services to be running
    # They are included in integration tests instead
    print_success "Functional tests completed (API Gateway tests skipped - require full system)"
}

# Function to run integration tests
run_integration_tests() {
    print_status "Running integration tests..."
    print_warning "Integration tests require the full system to be running"
    
    # Test API Gateway integration
    $TEST_CMD -tags=integration ./services/api-gateway/...
    if [ $? -eq 0 ]; then
        print_success "API Gateway integration tests passed"
    else
        print_error "API Gateway integration tests failed"
        return 1
    fi
    
    # Test repository integration
    $TEST_CMD -tags=integration ./services/user-service/internal/infrastructure/repository/...
    if [ $? -eq 0 ]; then
        print_success "Repository integration tests passed"
    else
        print_error "Repository integration tests failed"
        return 1
    fi
}

# Function to generate coverage report
generate_coverage_report() {
    if [ "$COVERAGE" = true ]; then
        print_status "Generating coverage report..."
        go tool cover -html=coverage.out -o coverage.html
        print_success "Coverage report generated: coverage.html"
    fi
}

# Main execution
print_status "Starting test execution in mode: $MODE"

case $MODE in
    unit)
        run_unit_tests
        ;;
    functional)
        run_functional_tests
        ;;
    integration)
        run_integration_tests
        ;;
    all)
        run_unit_tests
        run_functional_tests
        if [ "$INTEGRATION" = true ]; then
            run_integration_tests
        fi
        ;;
esac

# Generate coverage report if requested
generate_coverage_report

# Check if all tests passed
if [ $? -eq 0 ]; then
    print_success "All tests completed successfully!"
    exit 0
else
    print_error "Some tests failed!"
    exit 1
fi
