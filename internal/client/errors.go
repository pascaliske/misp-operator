package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ApiError represents a non-2xx response from the external API.
type ApiError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("api error (status %d): %s", e.StatusCode, e.Message)
}

func parseErrorResponse(res *http.Response) error {
	var apiErr ApiError

	// try to decode response body
	if err := json.NewDecoder(res.Body).Decode(&apiErr); err != nil {
		return &ApiError{
			StatusCode: res.StatusCode,
			Message:    res.Status,
		}
	}

	// inject status code
	apiErr.StatusCode = res.StatusCode

	return &apiErr
}
