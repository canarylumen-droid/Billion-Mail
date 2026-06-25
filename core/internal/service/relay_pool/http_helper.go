package relay_pool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

// doPost sends a JSON POST request and returns the response body + status code.
func doPost(ctx context.Context, url string, headers map[string]string, body interface{}) ([]byte, int, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return nil, 0, fmt.Errorf("encode body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, &buf)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	return respBody, resp.StatusCode, nil
}

// doPostForm sends an application/x-www-form-urlencoded POST.
func doPostForm(ctx context.Context, url string, headers map[string]string, formBody string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBufferString(formBody))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	return respBody, resp.StatusCode, nil
}

func errFromStatus(provider string, status int, body []byte) error {
	return fmt.Errorf("%s API error %d: %s", provider, status, truncate(string(body), 200))
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
