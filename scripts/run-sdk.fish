#!/usr/bin/fish

#   GO
#   $  ./scripts/generate-clean-sdk.sh --typescript --server
#   modules imports  *github.com/sst/opencode-sdk-go*   TO    *git.j9xym.com/opencode-api-go*
#   in *go.mod*, redirect   Module    TO    Directory 
#
#   rm sdk/go/types.gen.go    HAVING TO USE THE CURRENTLY EXISITING TOOLS
#
#
#
#   TYPESCRIPT
#   $  ./scripts/generate-clean-sdk.sh --typescript --server
#
#   
#   CLIENT-SIDE  TUI.EXE
#   $  go build -o packages/tui/opencode packages/tui/cmd/opencode/main.go
#   main.go  is the custom client application
#
#
#
##
# ├── sdk/go/                    # Fresh Go SDK (independent)
# │   ├── client.go             # Main client with services
# │   ├── types.gen.go          # Generated from OpenAPI
# │   ├── app.go, session.go, etc. # Service implementations
# │   └── integration_test.go   # End-to-end validation
# │
# ├── packages/tui/             # Main TUI application
# │   ├── cmd/opencode/main.go  # ✅ BUILT SUCCESSFULLY
# │   └── internal/             # Local TUI packages
# │       ├── app/, clipboard/, commands/, etc.
