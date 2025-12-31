package supabase

import (
	"context"
	"testing"
)

func TestCallDemoFunction(t *testing.T) {
	ctx := context.Background()
	supaClient := NewSupaClient(ctx)
	resp, err := supaClient.AddProjectRun(ctx, "0e586558-ba2e-4c3f-9129-08d5c61640b7", 4)
	if err != nil {
		t.Fatalf("Failed to add project run: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Failed to add project run: %v", resp.StatusCode)
	}

}
