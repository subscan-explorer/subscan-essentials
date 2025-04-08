package util

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

func PostWithJson(data []byte, endpoint string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.Body == nil {
		return nil, nil
	}
	defer resp.Body.Close() // nolint: errcheck
	body, _ := io.ReadAll(resp.Body)
	return body, nil
}

func HttpGet(ctx context.Context, endpoint string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.Body == nil {
		return nil, nil
	}
	defer resp.Body.Close() // nolint: errcheck
	if resp.StatusCode > 300 {
		return nil, fmt.Errorf("http error: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	return body, err
}
