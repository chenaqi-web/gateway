package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"backend/gateway/internal/config"
	"backend/gateway/internal/model/reponse"
	"backend/gateway/internal/model/request"
)

type PyClient struct {
	baseURL        string
	httpClient     *http.Client
	requestTimeout time.Duration
}

func NewHTTPClient(cfg *config.Config) *PyClient {
	timeoutSec := cfg.HTTP.RequestTimeout
	if timeoutSec <= 0 {
		timeoutSec = 30
	}

	baseURL := strings.TrimRight(cfg.HTTP.AgentServerAddr, "/")
	if baseURL == "" {
		baseURL = "http://127.0.0.1:8080"
	}

	return &PyClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Second * time.Duration(timeoutSec),
		},
		requestTimeout: time.Second * time.Duration(timeoutSec),
	}
}

// =====================================================================================================================

func (c *PyClient) GetRequestTimeout() time.Duration {
	return c.requestTimeout
}

func (c *PyClient) SearchVectors(
	ctx context.Context,
	collectionName string,
	req *request.VectorSearchRequest,
) (*reponse.VectorSearchResponse, error) {
	if collectionName == "" {
		return nil, fmt.Errorf("collection_name is required")
	}
	if req == nil || req.Content == "" {
		return nil, fmt.Errorf("content is required")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal search request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/vector/search/%s", c.baseURL, collectionName)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create search request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("call agent-server search: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read agent-server response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("agent-server status %d: %s", resp.StatusCode, string(respBody))
	}

	var result reponse.VectorSearchResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal agent-server response: %w", err)
	}
	if result.Code != 200 {
		return nil, fmt.Errorf("agent-server business error: code=%d message=%s", result.Code, result.Message)
	}

	return &result, nil
}
