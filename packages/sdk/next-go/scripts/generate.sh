#!/bin/bash
set -e

echo "Generating Go client..."

echo "Fetching OpenAPI spec..."
curl -o openapi.json http://127.0.0.1:3000/openapi.json

echo "Generating client..."
go generate ./...

echo "Done."
