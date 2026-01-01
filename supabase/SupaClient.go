package supabase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/mattmacf98/buildium_harness/logger"
)

type SupaClient struct {
	Client  *http.Client
	Token   string
	BaseUrl string
	AnonKey string
}

func NewSupaClient(ctx context.Context) *SupaClient {
	environment := os.Getenv("ENVIRONMENT")
	if environment == "PROD" {
		return &SupaClient{Client: http.DefaultClient, BaseUrl: "https://dpwumtpjesedslulexqz.supabase.co", AnonKey: os.Getenv("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImRwd3VtdHBqZXNlZHNsdWxleHF6Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjcxNDE4MzksImV4cCI6MjA4MjcxNzgzOX0.JYXW1bzTOmlCtngrlYLAbnGzRXDIcH0mDlwpbg1u8Rs")}
	} else {
		return &SupaClient{Client: http.DefaultClient, BaseUrl: "http://127.0.0.1:54321", AnonKey: os.Getenv("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZS1kZW1vIiwicm9sZSI6ImFub24iLCJleHAiOjE5ODM4MTI5OTZ9.CRXP1A7WOeoJeXxjNni43kdQwgnWNReilDMblYTn_I0")}
	}
}

func (c *SupaClient) AddProjectRun(ctx context.Context, projectId string, stage int, logs []logger.Log) (*http.Response, error) {
	logsJson, err := json.Marshal(logs)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseUrl+"/functions/v1/create-project-run",
		strings.NewReader(fmt.Sprintf(`{"projectId":"%s", "stage":%d, "logsJson":%s, "token":"%s"}`, projectId, stage, logsJson, c.Token)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-buildium-token", c.Token)
	req.Header.Set("Authorization", "Bearer "+c.AnonKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *SupaClient) Login(ctx context.Context, email string, password string) error {
	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseUrl+"/functions/v1/login",
		strings.NewReader(fmt.Sprintf(`{"email":"%s", "password":"%s"}`, email, password)))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.AnonKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to login: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var token struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return err
	}
	c.Token = token.Token
	return nil
}
