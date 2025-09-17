# OpenCode Go SDK Testing

This document describes the testing harness implemented for the OpenCode Go SDK.

## Test Files

We've created several test files to ensure comprehensive coverage:

1. `opencode_test.go` - Unit tests for the main SDK functions
2. `integration_test.go` - Integration tests that run against a live or mock server
3. `coverage_test.go` - Tools to track API implementation and test coverage

## How the Testing Harness Works

### Unit Tests

Unit tests in `opencode_test.go` verify the functionality of individual SDK functions without requiring a server connection. These tests check:

- Client creation
- Parameter handling
- Return value processing
- Error handling

### Integration Tests

Integration tests in `integration_test.go` verify the SDK's functionality against a server. These tests can run against:

1. A live server (if available)
2. A mock server (always available)

The mock server is implemented using Go's `httptest` package and simulates the API endpoints that the SDK interacts with.

### Coverage Tracking

The coverage tracking system in `coverage_test.go` helps measure how much of the API has been implemented and tested. It tracks:

- Total endpoints in the API
- Implemented endpoints
- Tested endpoints
- Coverage percentages

## How to Run the Tests

### Running All Tests

To run all tests, execute the test script:

```bash
./scripts/test.sh
```

Or run directly with Go:

```bash
go test -v ./...
```

### Running Tests with Coverage

To run tests with coverage reporting:

```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Running Specific Tests

To run a specific test:

```bash
go test -v -run TestName ./...
```

## Coverage Tracking

The coverage tracking system provides insights into the SDK's implementation status:

1. **Implementation Coverage** - Measures which API endpoints have been implemented in the SDK
2. **Test Coverage** - Measures which implemented endpoints have been tested

Coverage reports are generated in both JSON and HTML formats:

- `coverage_report.json` - Detailed coverage data
- `coverage.html` - Visual coverage report

## Integration Patterns

### Mock Server Pattern

The testing harness uses Go's `httptest` package to create a mock server that simulates the API endpoints. This allows tests to run in isolation without requiring a live server.

### Live Server Compatibility

Tests are designed to work with a live server when available. The `testutil.CheckTestServer` function checks if a server is running and skips tests that require a server if it's not available.

### Composition Pattern

The SDK follows a composition pattern where the generated API client is wrapped by a higher-level client. This pattern allows for extending functionality without modifying the generated code.

## Test Organization

Tests follow Go's standard conventions:

1. **File Naming** - Test files end with `_test.go`
2. **Function Naming** - Test functions start with `Test`
3. **Table-Driven Tests** - Multiple test cases are organized in tables for better maintainability
4. **Subtests** - Related tests are grouped using subtests
5. **Error Reporting** - Clear error messages with expected vs. actual values

## Best Practices Implemented

1. **Isolated Tests** - Each test can run independently
2. **Deterministic** - Tests produce consistent results
3. **Fast** - Unit tests run quickly without network calls
4. **Comprehensive** - Both positive and negative test cases
5. **Clear Error Messages** - Tests provide helpful error messages
6. **Coverage Tracking** - Implementation and test coverage are measured
