# API Testing Framework

This directory contains tests for the generated Ogen client in `../ogen/`.

## Running Tests

```bash
cd packages/tui
go test ./api/ogen -v -cover
```

For full coverage:

```bash
go test ./api/ogen -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Structure

- `client_test.go`: Unit/integration tests for all client methods using httptest mock server.
- Spec-driven: Uses kin-openapi to enumerate operations and generate test skeletons.
- Coverage: Aims for 100% on client methods; validates requests, responses, errors.

## Mock Server

Uses `httptest.Server` with a simple mux. For realism, generate Ogen server handlers:

```bash
ogen --generate-server --target ../server --package server openapi.json
```

Then integrate handlers into tests.

## Recommendations

- CI: Add to .github/workflows/test.yml: `go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out | grep total`.
- Extend to E2E: Use dockertest to spin TS server for real interactions.
- Fuzzing: Integrate go-fuzz for schema-based input generation.

Coverage goal: 100% on generated code.
