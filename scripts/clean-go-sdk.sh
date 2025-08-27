#!/bin/bash
# Clean generated Go SDK artifacts
set -euo pipefail
ROOT=$(cd "$(dirname "$0")/.." && pwd)
cd "$ROOT"

echo "Removing generated files..."
rm -f sdk/go/types.gen.go
rm -rf sdk/go/generated

echo "Removing module cache..."
go clean -modcache || true

if compgen -G "sdk/go/.go-cleanup-backup-*" > /dev/null; then
  echo "Backup exists; not restoring automatically. Use manual restore if needed."
fi

echo "Syncing workspace if go.work exists..."
if [ -f go.work ]; then
  go work sync
fi

echo "Clean complete"
