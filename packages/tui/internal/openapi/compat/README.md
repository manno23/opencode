# Compatibility layer

The compatibility layer provides adaptors and helpers for seamless integration between the hybrid Go SDK and existing TUI applications.

---

## Purpose

The opencode project uses a hybrid SDK generation approach that combines:

- **Auto-generated types** from OpenAPI specifications
- **Manual service implementations** for optimal developer experience

This creates a mismatch with existing TUI code that expects different type patterns. The compatibility layer bridges this gap without requiring extensive TUI refactoring.

---

## Architecture

### Type adaptors

Convert between generated SDK types and TUI-compatible types:

```go
// Generated SDK type (uses pointers, complex enums)
type Config struct {
    Share *ConfigShare `json:"share,omitempty"`
}

// TUI-compatible type (uses values, simple strings)
type CompatibleConfig struct {
    Share ConfigShare `json:"share,omitempty"`
}
```

### Method wrappers

Provide familiar method signatures for TUI integration:

```go
// Direct SDK usage (complex)
params := opencode.AppLogJSONRequestBody{
    Level:   opencode.AppLogJSONBodyLevel(level),
    Message: message,
    Service: service,
}
response, err := client.AppLogWithResponse(ctx, params)

// Compatibility layer (simple)
compatParams := compat.AppLogParams{
    Level:   compat.AppLogParamsLevelInfo,
    Message: "Hello world",
    Service: "tui",
}
err := compat.LogMessage(ctx, client, compatParams)
```

---

## Usage examples

### Config conversion

```go
import "github.com/sst/opencode-api-go/compat"

// Convert from generated SDK to TUI-compatible
func handleGeneratedConfig(sdkConfig *opencode.Config) {
    compatConfig := compat.ConvertCompatibleConfigFromGenerated(sdkConfig)

    // Now use simple value-based enum
    if compatConfig.Share == compat.ConfigShareDisabled {
        // Handle disabled sharing
    }
}
```

### Logging integration

```go
import "github.com/sst/opencode-api-go/compat"

// TUI logging with compatibility layer
func logToAPI(client *opencode.Client, level string, message string) {
    params := compat.AppLogParams{
        Level:   compat.AppLogParamsLevel(level),
        Message: message,
        Service: "tui",
        Extra:   map[string]interface{}{
            "component": "chat",
        },
    }

    // Convert and send to API
    sdkParams := compat.ConvertAppLogParamsToGenerated(params)
    client.App.Log(context.TODO(), sdkParams)
}
```

### Client initialization

```go
import "github.com/sst/opencode-api-go/compat"

// Simplified client creation
func initAPIClient(baseURL string) *opencode.Client {
    return compat.NewClient(baseURL)
}

// Session operations
func createSession(client *opencode.Client) error {
    params := map[string]any{
        "provider": "anthropic",
        "model":    "claude-3-sonnet",
    }

    _, err := compat.SessionInit(context.TODO(), client, params)
    return err
}
```

---

## Migration patterns

### From old SDK

**Before:**

```go
import "github.com/sst/opencode-api-go"

client := opencode.NewClient(url)
sessions, err := client.Session.List(context.TODO())
```

**After:**

```go
import "github.com/sst/opencode-api-go"
import "github.com/sst/opencode-api-go/compat"

client := compat.NewClient(url)  // Compatibility wrapper
sessions, err := client.Session.List(context.TODO()) // Direct SDK call
```

### Type handling

**Before:**

```go
// Direct string comparison
if config.Share == "disabled" {
    // Handle disabled
}
```

**After:**

```go
// Use compatibility constants
compatConfig := compat.ConvertCompatibleConfigFromGenerated(config)
if compatConfig.Share == compat.ConfigShareDisabled {
    // Handle disabled
}
```

---

## Available adaptors

### Config types

```go
// Enum mappings
const (
    ConfigShareAuto     ConfigShare = "auto"
    ConfigShareDisabled ConfigShare = "disabled"
    ConfigShareManual   ConfigShare = "manual"
)

// Conversion functions
func ConvertCompatibleConfigFromGenerated(*opencode.Config) CompatibleConfig
```

### Logging types

```go
// Log level mappings
const (
    AppLogParamsLevelDebug AppLogParamsLevel = "debug"
    AppLogParamsLevelInfo  AppLogParamsLevel = "info"
    AppLogParamsLevelError AppLogParamsLevel = "error"
    AppLogParamsLevelWarn  AppLogParamsLevel = "warn"
)

// Conversion functions
func ConvertAppLogParamsToGenerated(AppLogParams) opencode.AppLogJSONRequestBody
```

### Client helpers

```go
// Simplified constructors
func NewClient(baseURL string) *opencode.Client

// Common operations
func SessionInit(ctx context.Context, c *opencode.Client, params map[string]any) (any, error)
```

---

## Implementation guidelines

### Adding new adaptors

1. **Identify type mismatches** between generated SDK and TUI expectations
2. **Create value-based equivalent** for complex generated types
3. **Implement conversion functions** in both directions if needed
4. **Add constants** for enum-like values
5. **Write usage examples** and update documentation

### Type design principles

- **Use values over pointers** where possible for TUI compatibility
- **Use string constants** instead of complex enum types
- **Keep interfaces simple** - single responsibility per function
- **Maintain backward compatibility** with existing TUI code

### Error handling

```go
// Compatibility layer should handle SDK errors gracefully
func SafeConvert(sdkType *opencode.ComplexType) (CompatType, error) {
    if sdkType == nil {
        return CompatType{}, nil // Safe default
    }

    // Convert and validate
    result := CompatType{
        Field: sdkType.Field,
    }

    return result, nil
}
```

---

## Testing

### Unit tests

```go
func TestConfigConversion(t *testing.T) {
    sdkConfig := &opencode.Config{
        Share: &opencode.ConfigShareDisabled,
    }

    compat := ConvertCompatibleConfigFromGenerated(sdkConfig)

    assert.Equal(t, compat.ConfigShareDisabled, compat.Share)
}
```

### Integration tests

```go
func TestTUIIntegration(t *testing.T) {
    client := compat.NewClient("http://localhost:8080")

    // Test that TUI operations work through compatibility layer
    params := compat.AppLogParams{
        Level:   compat.AppLogParamsLevelInfo,
        Message: "Test message",
        Service: "test",
    }

    err := compat.LogMessage(context.TODO(), client, params)
    assert.NoError(t, err)
}
```

---

## Troubleshooting

### Common issues

**Type conversion errors:**

```go
// Problem: Generated type has unexpected structure
// Solution: Check OpenAPI spec and update compatibility layer

// Example fix:
func HandleTypeEvolution(newSDKType *opencode.NewType) CompatType {
    // Handle both old and new SDK patterns
    if newSDKType.NewField != nil {
        return CompatType{Field: *newSDKType.NewField}
    }
    return CompatType{Field: newSDKType.LegacyField}
}
```

**Missing constants:**

```go
// Problem: TUI expects constants not in compatibility layer
// Solution: Add constants to compat package

const (
    NewConstantValue CompatType = "new_value"
)
```

**Performance issues:**

```go
// Problem: Too many type conversions in hot paths
// Solution: Cache converted values or use direct SDK types

var configCache = make(map[string]CompatibleConfig)

func GetCompatConfig(key string, sdkConfig *opencode.Config) CompatibleConfig {
    if cached, ok := configCache[key]; ok {
        return cached
    }

    compat := ConvertCompatibleConfigFromGenerated(sdkConfig)
    configCache[key] = compat
    return compat
}
```

### Debug helpers

```go
// Add debug logging to understand type conversions
func debugConversion(from interface{}, to interface{}) {
    log.Printf("Converting %T -> %T: %+v -> %+v", from, to, from, to)
}
```
