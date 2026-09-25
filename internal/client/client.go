package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type ApiClient struct {
	baseUrl    *url.URL
	httpClient *http.Client
	apiKey     string
}

func NewApiClient(baseUrl string, apiKey string) (*ApiClient, error) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	parsedUrl, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	return &ApiClient{
		baseUrl:    parsedUrl,
		httpClient: httpClient,
		apiKey:     apiKey,
	}, nil
}

func (c *ApiClient) GetBaseUrl() *url.URL {
	return c.baseUrl
}

func (c *ApiClient) do(ctx context.Context, method string, endpoint string, body any, result any) error {
	// prepare request body
	reqBody, err := prepareRequestBody(body)
	if err != nil {
		return fmt.Errorf("encoding request body: %w", err)
	}

	// create request object
	req, err := http.NewRequestWithContext(ctx, method, c.baseUrl.JoinPath(endpoint).String(), reqBody)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	// setup request headers
	req.Header.Set("Authorization", c.apiKey)
	req.Header.Set("User-Agent", "misp-operator/1.0")
	req.Header.Set("Accept", "application/json")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// perform request
	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}

	// clean-up response body
	defer res.Body.Close()

	// handle error responses
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return parseErrorResponse(res)
	}

	// decode json response body
	if err := json.NewDecoder(res.Body).Decode(result); err != nil {
		return fmt.Errorf("decoding response body: %w", err)
	}

	return nil
}

func prepareRequestBody(body any) (io.Reader, error) {
	// no body, return empty value
	if body == nil {
		return nil, nil
	}

	buf := &bytes.Buffer{}

	// json encode body
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, fmt.Errorf("encoding request body: %w", err)
	}

	return buf, nil
}
