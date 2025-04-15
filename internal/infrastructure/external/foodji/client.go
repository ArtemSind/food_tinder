package foodji

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// DefaultBaseURL is the base URL for the Foodji API
	DefaultBaseURL = "https://amperoid.tenants.foodji.io"
	// DefaultTimeout is the default timeout for API requests
	DefaultTimeout = 10 * time.Second
)

// Client is a client for the Foodji API
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new Foodji API client
func NewClient() *Client {
	return &Client{
		BaseURL: DefaultBaseURL,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// GetMachine fetches machine data by ID
func (c *Client) GetMachineProducts(ctx context.Context, machineID string) (*[]Product, error) {
	url := fmt.Sprintf("%s/machines/%s", c.BaseURL, machineID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var apiResponse Response
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &apiResponse.Data.MachineProducts, nil
}
