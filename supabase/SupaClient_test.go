package supabase

import (
	"context"
	"testing"
)

func TestCallDemoFunction(t *testing.T) {
	ctx := context.Background()
	supaClient := NewSupaClient(ctx)
	resp, err := supaClient.CallDemoFunction(ctx, "TEST")
	if err != nil {
		t.Fatalf("Failed to call demo function: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Failed to call demo function: %v", resp.StatusCode)
	}

}
