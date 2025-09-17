#!/bin/bash
set -e

echo "Running OpenCode Go SDK tests..."

# Run unit tests
echo "Running unit tests..."
go test -v ./test/...

# Run coverage tests
echo "Running coverage tests..."
go test -v -coverprofile=coverage.out ./test/...

# Generate coverage report
echo "Generating coverage report..."
go tool cover -html=coverage.out -o coverage.html

echo "Tests completed successfully!"
echo "Coverage report saved to coverage.html"