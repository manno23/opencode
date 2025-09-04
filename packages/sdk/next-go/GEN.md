# Unique Benefits of Ogen Over Other Go Generators



```.ogen.yml
generator:
  features:
    disable_all: true  # Start fresh, enable selectively
    enable:
      - "paths/client"               # Generate client code
      - "paths/server"               # Generate server code
      - "client/request/validation"  # Enable request validation
      - "server/response/validation" # Enable response validation
      - "ogen/otel"                  # OpenTelemetry integration
  infer_types: true                  # Auto-infer schema types
```






```
// Custom adapter wrapping the generated client.

package adapters

import (
	"context"
	"log"

	"yourproject/api"
)

type ClientAdapter struct {
	client *api.Client
}

func NewClientAdapter(baseURL string) (*ClientAdapter, error) {
	cl, err := api.NewClient(baseURL, api.WithHTTPClient(http.DefaultClient))
	if err != nil {
		return nil, err
	}
	return &ClientAdapter{client: cl}, nil
}

func (a *ClientAdapter) GetPetWithLogging(ctx context.Context, id int64) (api.CustomPet, error) {
	log.Printf("Fetching pet with ID: %d", id)
	params := api.GetPetParams{ID: id}
	pet, err := a.client.GetPet(ctx, params)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	return pet, err
}

// Extend for auth, retries, etc.
```


- adapters/server_handler_impl.go
```
package adapters

import (
	"context"
	"io"
	"net/http"
	"yourproject/api"
)

type MyHandler struct {
	// Dependencies like DB...
}

func (h *MyHandler) GetPet(ctx context.Context, params api.GetPetParams) (api.CustomPet, error) {
	// Business logic: Fetch from DB...
	return api.CustomPet{Name: "Fluffy", Age: 5}, nil
}

func (h *MyHandler) StreamEvents(ctx context.Context) (io.ReadCloser, error) {
	// Server-side SSE: Create a pipe and write events
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		// Simulate events
		pw.Write([]byte("data: Event 1\n\n"))
		pw.Write([]byte("data: Event 2\n\n"))
	}()
	return pr, nil
}

// Usage: http.ListenAndServe(":8080", api.NewServer(&MyHandler{}))
```

## Ogen distinguishes itself with:

- **No Reflection or Dynamic Typing**: Pure code-generated JSON marshaling/unmarshaling for speeds up to GB/s, unlike oapi-codegen which relies on reflection leading to runtime overhead.
- **Static Radix Router**: Efficient routing for servers, outperforming dynamic routers in tools like OpenAPI Generator.
- **Sum Types for oneOf/anyOf**: Automatic type discrimination with tagged unions, reducing manual checks and errors.
- **Generic Wrappers for Optionals/Nullables**: Avoids pointer indirection pitfalls, minimizing GC pressure.
- **Built-in Validation and Error Handling**: Detailed, type-safe errors without extra libs.
- **Streaming Support**: Native for SSE/JSON streams, with low memory usage for large responses.
- **Extensibility**: Features like OpenTelemetry and custom extensions make it ideal for production.

Compared to alternatives:
- **oapi-codegen**: Simpler but slower due to reflection; lacks sum types.
- **OpenAPI Generator**: Broad language support but generates verbose, less performant Go code.
- **Trade-offs**: Ogen may require more spec annotations for full optimization, but yields superior runtime efficiency for high-throughput apps.

This makes ogen ideal for type-strict, performant APIs with real-time features like SSE.
