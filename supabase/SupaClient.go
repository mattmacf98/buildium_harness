package supabase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mattmacf98/buildium_harness/logger"
)

type SupaClient struct {
	Client *http.Client
}

func NewSupaClient(ctx context.Context) *SupaClient {
	return &SupaClient{Client: http.DefaultClient}
}

func (c *SupaClient) AddProjectRun(ctx context.Context, projectId string, stage int, logs []logger.Log) (*http.Response, error) {
	logsJson, err := json.Marshal(logs)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", "https://dpwumtpjesedslulexqz.supabase.co/functions/v1/create-project-run",
		strings.NewReader(fmt.Sprintf(`{"projectId":"%s", "stage":%d, "logsJson":%s}`, projectId, stage, logsJson)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImRwd3VtdHBqZXNlZHNsdWxleHF6Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjcxNDE4MzksImV4cCI6MjA4MjcxNzgzOX0.JYXW1bzTOmlCtngrlYLAbnGzRXDIcH0mDlwpbg1u8Rs")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
