# OpenCode Go SDK Tests

This directory contains all the tests for the OpenCode Go SDK.

## Test Organization

The tests are organized as follows:

1. `opencode_test.go` - Unit tests for the main SDK functionality
2. `integration_test.go` - Integration tests that can run against a live or mock server
3. `coverage_test.go` - Coverage tracking mechanisms and reporting

## Running Tests

To run all tests:

```bash
cd /home/jm/data/code/opencode/packages/sdk/next-go
./scripts/test.sh
```

Or run directly with Go:

```bash
go test -v ./test/...
```

## Test Structure

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

## Test Patterns Used

1. **Mock Server Pattern** - Uses `httptest` for isolated testing without external dependencies
2. **Live Server Compatibility** - Tests can work with live servers when available
3. **Table-Driven Tests** - Multiple test cases organized in tables for better maintainability
4. **Subtests** - Related tests grouped using subtests
5. **Coverage Tracking** - Implementation and test coverage monitoring
