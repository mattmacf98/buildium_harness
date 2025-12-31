package supabase

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type SupaClient struct {
	Client *http.Client
}

func NewSupaClient(ctx context.Context) *SupaClient {
	return &SupaClient{Client: http.DefaultClient}
}

func (c *SupaClient) CallDemoFunction(ctx context.Context, name string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", "https://dpwumtpjesedslulexqz.supabase.co/functions/v1/hello-world",
		strings.NewReader(fmt.Sprintf(`{"name":"%s"}`, name)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer sb_publishable_KJy0pE3QshfkUgn_QHrVYQ_BoiLh2i0")
	req.Header.Set("apikey", "sb_publishable_KJy0pE3QshfkUgn_QHrVYQ_BoiLh2i0")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil

}
