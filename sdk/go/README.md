# OpenCode Go SDK

The OpenCode Go SDK provides convenient access to the [OpenCode REST API](https://opencode.ai/docs) from applications written in Go.

This SDK is generated using a clean, independent approach that ensures complete API coverage and modern Go patterns.

## Installation

```go
import (
    "git.j9xym.com/opencode-api-go" // imported as opencode
)
```

## Requirements

This library requires Go 1.24+.

## Usage

```go
package main

import (
    "context"
    "fmt"

    "git.j9xym.com/opencode-api-go"
)

func main() {
    client := opencode.NewClient()

    // Get app information
    app, err := client.App.Get(context.Background())
    if err != nil {
        panic(err)
    }
    fmt.Printf("App: %+v\n", app)

    // Create a new session
    session, err := client.Session.New(context.Background(), opencode.SessionNewParams{})
    if err != nil {
        panic(err)
    }
    fmt.Printf("Session: %+v\n", session)
}
```

## Configuration

The client can be configured with various options:

```go
client := opencode.NewClient(
    opencode.WithBaseURL("https://api.opencode.dev"),
    // Add other options as needed
)
```

## Services

The SDK provides the following services:

- `App` - Application information and configuration
- `Session` - Session management
- `Command` - Command execution
- `File` - File operations
- `Find` - Search functionality
- `Config` - Configuration management
- `Event` - Server-sent events
- `Tui` - Terminal UI operations

## Error Handling

All methods return errors that should be handled appropriately:

```go
session, err := client.Session.New(context.Background(), params)
if err != nil {
    // Handle error
    return err
}
```

## Testing

The SDK includes comprehensive integration tests. To run them:

```bash
go test ./...
```

For integration tests against a live server:

```bash
go test -tags=integration ./integration_test.go
```

## Generation

This SDK is generated using the clean SDK generation script:

```bash
./scripts/generate-clean-sdk.sh --go --server
```

This ensures the SDK is always up-to-date with the latest API specification.
