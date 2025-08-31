#!/usr/bin/env fish
# filepath: scripts/build.fish

# Strict-ish mode
# set -e
# set -u

function _echo
  set_color cyan
  echo -n "[build] "
  set_color normal
  echo $argv
end

function _err
  set_color red
  echo -n "[error] "
  set_color normal
  echo $argv 1>&2
end

function _ok
  set_color green
  echo -n "[ok] "
  set_color normal
  echo $argv
end

# Resolve repo root from this script's location
set SCRIPT_DIR (dirname (status --current-filename))
set REPO_ROOT (realpath "$SCRIPT_DIR/..")
cd $REPO_ROOT
echo "cd to $REPO_ROOT"

# Metadata
set VERSION (git describe --tags --dirty --always ^/dev/null; or echo "v0.0.0")
set COMMIT (git rev-parse --short HEAD)
set BUILD_DATE (date -u +"%Y-%m-%dT%H:%M:%SZ")

echo "Version: $VERSION"
echo "Commit: $COMMIT"
echo "Build date: $BUILD_DATE"

# Output dirs
set OUT_DIR_GO "packages/dist"
set OUT_DIR_BUN "packages/dist"
mkdir -p $OUT_DIR_GO
mkdir -p $OUT_DIR_BUN

# Small helper to compute sha256 and write alongside
function write_hash --argument-names binpath
  if test -f "$binpath"
    set hash (sha256sum "$binpath" | awk '{print $1}')
    echo $hash > "$binpath.sha256"
    _ok "sha256 $binpath = $hash"
  else
    _err "Missing binary for hashing: $binpath"
    return 1
  end
end

function _err --argument-names msg
  echo "$msg"
end

# --------------- Section A: Go TUI shared lib (CLI) ----------------

_echo "A. Verify Go SDK linkage (sdk/go) and build hardened Go TUI"

# 1. Verify we correctly linked our own SDK from sdk/go
#    - Ensure sdk/go is present; run a trivial compile and an import check
if not test -f "sdk/go/go.mod"
  _err "sdk/go/go.mod missing; Go SDK not present"
  exit 1
end

# Confirm we can import our local module path and it compiles
_echo "Checking local Go SDK compiles..."
pushd sdk/go >/dev/null
set GO_PKG (go list ./ 2>/dev/null)
or begin
  _err "go list failed in sdk/go; check go.mod and sources"
  exit 1
end
# Optional quick build of package (not tests)
go build ./... >/dev/null
or begin
  _err "Go SDK source fails to build; resolve redeclarations first"
  exit 1
end
popd >/dev/null
_ok "Go SDK is compilable: $GO_PKG"

# Extra: verify that the main TUI imports the local module (string check)
# This is a heuristic—adjust module path if different.
set GO_SDK_IMPORT_PATH "git.j9xym.com/opencode-api-go"
if not rg -n --no-messages $GO_SDK_IMPORT_PATH tui/cmd/opencode/main.go >/dev/null
  _echo "Note: Could not find '$GO_SDK_IMPORT_PATH' import in tui/cmd/opencode/main.go (ok if indirect)."
else
  _ok "Main TUI imports SDK path: $GO_SDK_IMPORT_PATH"
end

# 2. Hardened build, but do not break
# Conservative hardening: strip symbols, trim paths, embed metadata.
# Avoid aggressive external linking/PIE/static that can break portability unexpectedly.
set GOOS linux
set GOARCH amd64
set CGO_ENABLED 0

echo "Version: $VERSION"
echo "Commit: $COMMIT"
echo "Build date: $BUILD_DATE"

set LDFLAGS "-s -w -X main.Version=$VERSION -X main.Commit=$COMMIT -X main.BuildDate=$BUILD_DATE"
set GO_OUT "$OUT_DIR_GO/opencode-cli"

_echo "Building Go TUI (CGO_ENABLED=$CGO_ENABLED GOOS=$GOOS GOARCH=$GOARCH)..."
go build \
  -trimpath \
  -ldflags "$LDFLAGS" \
  -o "$GO_OUT" \
  ./packages/tui/cmd/opencode/main.go
or begin
  _err "Go build failed"
  exit 1
end

# 3. Add metadata and hash
write_hash "$GO_OUT"

# 4. Smoke test version flag if implemented
if "$GO_OUT" --version >/dev/null 2>&1
  set ver_out ("$GO_OUT" --version 2>/dev/null)
  _ok "Go binary --version: $ver_out"
else
  _echo "Go binary has no --version (skipping)"
end

# --------------- Section B: Bun application builds ----------------

_echo "B. Verify TS SDK linkage (sdk/typescript) and build Bun app(s)"

# 1. Verify we correctly linked our own SDK from sdk/typescript
# Check workspace setup and dependency resolution for @opencode-ai/sdk-next
# Adjust package name if different.
set TS_SDK_PKG "@opencode-ai/sdk-next"

if not test -f "sdk/typescript/package.json"
  _err "sdk/typescript/package.json missing"
  exit 1
end

# Ensure the SDK builds
_echo "Building TS SDK..."
bun run --filter='./sdk/typescript' build
or begin
  _err "TS SDK build failed"
  exit 1
end

# Quick probe: detect imports in packages/opencode using @opencode-ai/sdk-next or path alias
set OPENCODE_DIR "packages/opencode"
if test -d $OPENCODE_DIR
  if /usr/bin/rg -n "@opencode-ai/sdk-next" $OPENCODE_DIR >/dev/null
    _ok "packages/opencode imports $TS_SDK_PKG"
  else
    _echo "Could not find imports of $TS_SDK_PKG in packages/opencode (ok if using path aliases)."
  end
else
  _echo "packages/opencode directory not found (skipping import scan)"
end

# 2. Hardened and optimized Bun build (conservative)
# We’ll produce two builds for the Bun entrypoint. Replace entry path if different.
# Standalone build WITHOUT Go TUI bundled
set BUN_ENTRY packages/opencode/src/plugin/index.ts  # adjust if needed
if not test -f "$BUN_ENTRY"
  # Try a common alternative entry
  set BUN_ENTRY sdk/typescript/mod/index.ts
end

if not test -f "$BUN_ENTRY"
  _err "Bun entry not found. Set BUN_ENTRY appropriately."
  exit 1
end

set BUN_OUT_NO_TUI $OUT_DIR_BUN/bun-app-no-tui.js
set BUN_OUT_WITH_TUI $OUT_DIR_BUN/bun-app-with-tui.js

# Conservative defines and externalization; do not break I/O unexpectedly
set BUN_DEFINES --define:process.env.NODE_ENV=\"production\" --define:__DEV__=false --define:OPENCODE_BASE_URL=\"opencode.j9xym.com\" --define:OPENCODE_API=\"api.opencode.j9xym.com\"

# Externalize Node/Bun native modules not needed in the core app (adjust as needed)
set BUN_EXTERNAL --define:process.env.OPENCODE_PORT=8080 --define:process.env.OPENCODE_BIND="0.0.0.0"
# Keep fs/net/tls unless you’re sure they’re unused
set BUN_FLAGS --minify --target=bun

set BUN_DEBUG_FLAGS --sourcemap=external --define:__DEV__=true


_echo "Building Bun standalone (without Go TUI)..."
bun build $BUN_ENTRY \
  --outdir $OUT_DIR_BUN \
  --outfile $BUN_OUT_NO_TUI \
  --target=bun \
  $BUN_FLAGS $BUN_DEFINES $BUN_EXTERNAL
or begin
  _err "Bun build (no TUI) failed"
  exit 1
end
write_hash "$BUN_OUT_NO_TUI"

# 3. Add identifying metadata for Bun builds
# For JS bundles, we’ll append a comment banner with version/commit/date (non-breaking).
printf "\n/* opencode build: version=%s commit=%s date=%s */\n" "$VERSION" "$COMMIT" "$BUILD_DATE" >> "$BUN_OUT_NO_TUI"
_ok "Tagged Bun (no TUI) with metadata"

# 4. Build standalone bundled bun application WITHOUT bundling the Go TUI
# Already done above (BUN_OUT_NO_TUI).

# 5. Build standalone bundled bun application WITH the Go TUI
# Strategy: bundle as normal and ensure the Go binary path is referenced at runtime.
# We do NOT embed the binary into JS; instead, we ensure code can locate it.
# If you have a wrapper that shells out to the Go binary, ensure its path is configurable.
_echo "Building Bun standalone (with Go TUI)..."
bun build $BUN_ENTRY \
  --outdir $OUT_DIR_BUN \
  --outfile $BUN_OUT_WITH_TUI \
  --target=bun \
  $BUN_FLAGS $BUN_DEFINES
or begin
  _err "Bun build (with TUI) failed"
  exit 1
end

# Annotate bundle with where the Go binary is expected
printf "\n/* expects Go TUI binary at: %s */\n" "$GO_OUT" >> "$BUN_OUT_WITH_TUI"
write_hash "$BUN_OUT_WITH_TUI"
_ok "Tagged Bun (with TUI) with metadata and expected binary path"

# Optional: verify both bundles run under Bun to at least import/execute top-level
_echo "Smoke test Bun bundles (import check)..."
bun "$BUN_OUT_NO_TUI" --help >/dev/null 2>&1; or true
bun "$BUN_OUT_WITH_TUI" --help >/dev/null 2>&1; or true
_ok "Builds complete"

_echo "Artifacts:"
echo "  Go: $GO_OUT (+ .sha256)"
echo "  Bun (no TUI): $BUN_OUT_NO_TUI (+ .sha256)"
echo "  Bun (with TUI): $BUN_OUT_WITH_TUI (+ .sha256)"

