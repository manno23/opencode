package test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// CoverageTracker tracks which API endpoints have been implemented and tested
type CoverageTracker struct {
	TotalEndpoints   int            `json:"total_endpoints"`
	Implemented      int            `json:"implemented"`
	Tested           int            `json:"tested"`
	EndpointCoverage map[string]int `json:"endpoint_coverage"` // 0 = not implemented, 1 = implemented, 2 = tested
}

// NewCoverageTracker creates a new coverage tracker
func NewCoverageTracker() *CoverageTracker {
	return &CoverageTracker{
		EndpointCoverage: make(map[string]int),
	}
}

// MarkImplemented marks an endpoint as implemented
func (c *CoverageTracker) MarkImplemented(endpoint string) {
	if c.EndpointCoverage[endpoint] < 1 {
		c.EndpointCoverage[endpoint] = 1
		c.Implemented++
	}
}

// MarkTested marks an endpoint as tested
func (c *CoverageTracker) MarkTested(endpoint string) {
	if c.EndpointCoverage[endpoint] < 2 {
		if c.EndpointCoverage[endpoint] == 1 {
			c.Tested++
		}
		c.EndpointCoverage[endpoint] = 2
	}
}

// CalculateCoverage calculates the overall coverage percentage
func (c *CoverageTracker) CalculateCoverage() float64 {
	if c.TotalEndpoints == 0 {
		return 0
	}
	return float64(c.Implemented) / float64(c.TotalEndpoints) * 100
}

// CalculateTestCoverage calculates the test coverage percentage
func (c *CoverageTracker) CalculateTestCoverage() float64 {
	if c.TotalEndpoints == 0 {
		return 0
	}
	return float64(c.Tested) / float64(c.TotalEndpoints) * 100
}

// Save saves the coverage report to a file
func (c *CoverageTracker) Save(filename string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// Load loads a coverage report from a file
func (c *CoverageTracker) Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, c)
}

// TestCoverageTracker tests the coverage tracker functionality
func TestCoverageTracker(t *testing.T) {
	tracker := NewCoverageTracker()
	tracker.TotalEndpoints = 10

	// Test marking endpoints as implemented
	tracker.MarkImplemented("/api/v1/app/log")
	tracker.MarkImplemented("/api/v1/session/get")
	tracker.MarkImplemented("/api/v1/session/permissions/respond")

	if tracker.Implemented != 3 {
		t.Errorf("Expected 3 implemented endpoints, got %d", tracker.Implemented)
	}

	// Test marking endpoints as tested
	tracker.MarkTested("/api/v1/app/log")
	tracker.MarkTested("/api/v1/session/get")

	if tracker.Tested != 2 {
		t.Errorf("Expected 2 tested endpoints, got %d", tracker.Tested)
	}

	// Test coverage calculation
	coverage := tracker.CalculateCoverage()
	expected := 30.0 // 3 out of 10
	if coverage != expected {
		t.Errorf("Expected coverage %.2f%%, got %.2f%%", expected, coverage)
	}

	testCoverage := tracker.CalculateTestCoverage()
	expectedTest := 20.0 // 2 out of 10
	if testCoverage != expectedTest {
		t.Errorf("Expected test coverage %.2f%%, got %.2f%%", expectedTest, testCoverage)
	}

	// Test saving and loading
	tempDir := t.TempDir()
	reportFile := filepath.Join(tempDir, "coverage.json")

	if err := tracker.Save(reportFile); err != nil {
		t.Fatalf("Failed to save coverage report: %v", err)
	}

	newTracker := NewCoverageTracker()
	if err := newTracker.Load(reportFile); err != nil {
		t.Fatalf("Failed to load coverage report: %v", err)
	}

	if newTracker.Implemented != tracker.Implemented {
		t.Errorf("Loaded tracker has different implemented count: expected %d, got %d",
			tracker.Implemented, newTracker.Implemented)
	}

	if newTracker.Tested != tracker.Tested {
		t.Errorf("Loaded tracker has different tested count: expected %d, got %d",
			tracker.Tested, newTracker.Tested)
	}
}
