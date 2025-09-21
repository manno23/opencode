package client

import (
	"context"
	"fmt"
	"io"

	"github.com/sst/opencode-sdk-go"
)

// FilesClient defines the interface for file operations
type FilesClient interface {
	Upload(ctx context.Context, reader io.Reader, meta FileMetadata) (*File, error)
	Download(ctx context.Context, fileID string) (io.ReadCloser, error)
	List(ctx context.Context, opts FileListOptions) ([]*File, error)
	Delete(ctx context.Context, fileID string) error
	Get(ctx context.Context, fileID string) (*File, error)
}

// File represents a file in the facade layer
type File struct {
	ID          string
	Name        string
	Size        int64
	ContentType string
	CreatedAt   string
	UpdatedAt   string
	URL         *string
}

// FileMetadata contains metadata for file uploads
type FileMetadata struct {
	Name        string
	ContentType string
	Size        *int64
}

// FileListOptions contains options for listing files
type FileListOptions struct {
	Limit  *int
	Offset *int
}

// LegacyFiles implements FilesClient using the legacy SDK
type LegacyFiles struct {
	client *opencode.Client
}

// NewLegacyFiles creates a new legacy files client
func NewLegacyFiles(client *opencode.Client) FilesClient {
	return &LegacyFiles{client: client}
}

// Upload is not supported by the legacy SDK. Present for interface compatibility.
func (l *LegacyFiles) Upload(ctx context.Context, reader io.Reader, meta FileMetadata) (*File, error) {
	return nil, fmt.Errorf("legacy SDK does not support file upload")
}

// Download is not supported by the legacy SDK. Present for interface compatibility.
func (l *LegacyFiles) Download(ctx context.Context, fileID string) (io.ReadCloser, error) {
	return nil, fmt.Errorf("legacy SDK does not support file download")
}

// List is not supported by the legacy SDK in this facade format. Present for interface compatibility.
func (l *LegacyFiles) List(ctx context.Context, opts FileListOptions) ([]*File, error) {
	return []*File{}, nil
}

// Delete is not supported by the legacy SDK. Present for interface compatibility.
func (l *LegacyFiles) Delete(ctx context.Context, fileID string) error {
	return fmt.Errorf("legacy SDK does not support file delete")
}

// Get is not supported by the legacy SDK. Present for interface compatibility.
func (l *LegacyFiles) Get(ctx context.Context, fileID string) (*File, error) {
	return nil, fmt.Errorf("legacy SDK does not support file get")
}
