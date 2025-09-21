package client

import (
	"context"
	"fmt"

	"github.com/sst/opencode-sdk-go"
)

// ToolsClient defines the interface for tool operations
type ToolsClient interface {
	List(ctx context.Context) ([]*Tool, error)
	Invoke(ctx context.Context, toolID string, params interface{}) (*ToolResult, error)
	GetPermissions(ctx context.Context, sessionID string) ([]*Permission, error)
	RespondToPermission(ctx context.Context, sessionID, permissionID string, response PermissionResponse) error
}

// Tool represents a tool in the facade layer
type Tool struct {
	ID          string
	Name        string
	Description string
	Schema      interface{}
}

// ToolResult represents the result of a tool invocation
type ToolResult struct {
	Success bool
	Output  interface{}
	Error   *string
}

// Permission represents a permission request
type Permission struct {
	ID          string
	ToolID      string
	SessionID   string
	Description string
	Status      PermissionStatus
}

// PermissionStatus represents permission states
type PermissionStatus string

const (
	PermissionStatusPending  PermissionStatus = "pending"
	PermissionStatusApproved PermissionStatus = "approved"
	PermissionStatusDenied   PermissionStatus = "denied"
)

// PermissionResponse represents a response to a permission request
type PermissionResponse string

const (
	PermissionResponseApprove PermissionResponse = "approve"
	PermissionResponseDeny    PermissionResponse = "deny"
)

// LegacyTools implements ToolsClient using the legacy SDK
type LegacyTools struct {
	client *opencode.Client
}

// NewLegacyTools creates a new legacy tools client
func NewLegacyTools(client *opencode.Client) ToolsClient {
	return &LegacyTools{client: client}
}

func (l *LegacyTools) List(ctx context.Context) ([]*Tool, error) {
	return []*Tool{}, nil
}

func (l *LegacyTools) Invoke(ctx context.Context, toolID string, params interface{}) (*ToolResult, error) {
	return nil, fmt.Errorf("legacy SDK does not support facade Tools.Invoke")
}

func (l *LegacyTools) GetPermissions(ctx context.Context, sessionID string) ([]*Permission, error) {
	return []*Permission{}, nil
}

func (l *LegacyTools) RespondToPermission(ctx context.Context, sessionID, permissionID string, response PermissionResponse) error {
	return fmt.Errorf("legacy SDK does not support facade Tools.RespondToPermission")
}
