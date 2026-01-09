package meta

import (
	"encoding/json"
	"testing"
)

func TestMetaStruct_JSONMarshalling(t *testing.T) {
	meta := Meta{
		Stage:         3,
		Entrypoint:    "main.go",
		ExecutableDir: "/some/path",
		ProjectId:     "proj-123",
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
	if result.ExecutableDir != meta.ExecutableDir {
		t.Errorf("Path = %q, want %q", result.ExecutableDir, meta.ExecutableDir)
	}
	if result.ProjectId != meta.ProjectId {
		t.Errorf("ProjectId = %q, want %q", result.ProjectId, meta.ProjectId)
	}
}

func TestMetaStruct_JSONTags(t *testing.T) {
	meta := Meta{
		Stage:         5,
		Entrypoint:    "app",
		ExecutableDir: "/test",
		ProjectId:     "abc-456",
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
