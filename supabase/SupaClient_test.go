package supabase

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/buildium-org/buildium_harness/logger"
)

func TestAddProjectRun(t *testing.T) {
	ctx := context.Background()
	supaClient := NewSupaClient(ctx)
	logs := []logger.Log{
		{Stage: 0, Message: "Health Check", Type: "HEADER"},
		{Stage: 1, Message: "test log", Type: "SUCCESS"},
		{Stage: 1, Message: "200 ok", Type: "SUCCESS"},
		{Stage: 2, Message: "Get Root", Type: "HEADER"},
		{Stage: 2, Message: "200 ok", Type: "SUCCESS"},
		{Stage: 2, Message: "Post Root", Type: "FAILURE"},
	}

	resp, err := supaClient.AddProjectRun(ctx, "0e586558-ba2e-4c3f-9129-08d5c61640b7", 4, logs)
	if err != nil {
		t.Fatalf("Failed to add project run: %v", err)
	}
	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}
		t.Log(string(body))
		t.Fatalf("Failed to add project run: %v", resp.StatusCode)

	}
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	supaClient := NewSupaClient(ctx)
	os.Setenv("BUILDIUM_EMAIL", "test@test.com")
	os.Setenv("BUILDIUM_PASSWORD", "test1234")
	err := supaClient.Login(ctx)
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	if supaClient.Token == "" {
		t.Fatalf("Failed to get token")
	}
}
