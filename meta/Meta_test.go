package meta

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestMetaStruct_JSONMarshalling(t *testing.T) {
	meta := Meta{
		Stage:      3,
		Entrypoint: "main.go",
		Path:       "/some/path",
		ProjectId:  "proj-123",
	}

	data, err := json.Marshal(meta)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var result Meta
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if result.Stage != meta.Stage {
		t.Errorf("Stage = %d, want %d", result.Stage, meta.Stage)
	}
	if result.Entrypoint != meta.Entrypoint {
		t.Errorf("Entrypoint = %q, want %q", result.Entrypoint, meta.Entrypoint)
	}
	if result.Path != meta.Path {
		t.Errorf("Path = %q, want %q", result.Path, meta.Path)
	}
	if result.ProjectId != meta.ProjectId {
		t.Errorf("ProjectId = %q, want %q", result.ProjectId, meta.ProjectId)
	}
}

func TestMetaStruct_JSONTags(t *testing.T) {
	meta := Meta{
		Stage:      5,
		Entrypoint: "app",
		Path:       "/test",
		ProjectId:  "abc-456",
	}

	data, err := json.Marshal(meta)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	// Verify JSON keys match the tags
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	expectedKeys := []string{"stage", "entrypoint", "path", "projectId"}
	for _, key := range expectedKeys {
		if _, ok := raw[key]; !ok {
			t.Errorf("expected JSON key %q not found", key)
		}
	}
}

func TestNewMeta(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "meta_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test meta.json file
	metaContent := `{
		"stage": 2,
		"entrypoint": "bin/app",
		"projectId": "test-project-789"
	}`

	metaPath := filepath.Join(tmpDir, "meta.json")
	if err := os.WriteFile(metaPath, []byte(metaContent), 0644); err != nil {
		t.Fatalf("failed to write meta.json: %v", err)
	}

	// Test NewMeta
	meta := NewMeta(tmpDir)

	if meta == nil {
		t.Fatal("NewMeta() returned nil")
	}

	if meta.Stage != 2 {
		t.Errorf("Stage = %d, want 2", meta.Stage)
	}
	if meta.Entrypoint != "bin/app" {
		t.Errorf("Entrypoint = %q, want %q", meta.Entrypoint, "bin/app")
	}
	if meta.ProjectId != "test-project-789" {
		t.Errorf("ProjectId = %q, want %q", meta.ProjectId, "test-project-789")
	}
	// Path should be set to the directory passed to NewMeta
	if meta.Path != tmpDir {
		t.Errorf("Path = %q, want %q", meta.Path, tmpDir)
	}
}

func TestNewMeta_SetsPathFromArgument(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "meta_test_path")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create meta.json with a different path value (should be overwritten)
	metaContent := `{
		"stage": 1,
		"entrypoint": "main",
		"path": "/original/path",
		"projectId": "proj-1"
	}`

	metaPath := filepath.Join(tmpDir, "meta.json")
	if err := os.WriteFile(metaPath, []byte(metaContent), 0644); err != nil {
		t.Fatalf("failed to write meta.json: %v", err)
	}

	meta := NewMeta(tmpDir)

	// Path should be overwritten with the argument, not the JSON value
	if meta.Path != tmpDir {
		t.Errorf("Path = %q, want %q (should override JSON value)", meta.Path, tmpDir)
	}
}

func TestNewMeta_MinimalJSON(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "meta_test_minimal")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create minimal meta.json (only required fields)
	metaContent := `{}`

	metaPath := filepath.Join(tmpDir, "meta.json")
	if err := os.WriteFile(metaPath, []byte(metaContent), 0644); err != nil {
		t.Fatalf("failed to write meta.json: %v", err)
	}

	meta := NewMeta(tmpDir)

	// Should have zero values for all fields except Path
	if meta.Stage != 0 {
		t.Errorf("Stage = %d, want 0 (zero value)", meta.Stage)
	}
	if meta.Entrypoint != "" {
		t.Errorf("Entrypoint = %q, want empty string", meta.Entrypoint)
	}
	if meta.ProjectId != "" {
		t.Errorf("ProjectId = %q, want empty string", meta.ProjectId)
	}
	if meta.Path != tmpDir {
		t.Errorf("Path = %q, want %q", meta.Path, tmpDir)
	}
}
