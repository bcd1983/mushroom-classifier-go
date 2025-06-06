// Package httpclient provides HTTP client utilities for making REST API calls
package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Request contains parameters for an HTTP POST request
//
// Contains all necessary information to make an authenticated
// HTTP POST request with a JSON payload.
type Request struct {
	// Target URL for the request
	URL string

	// Bearer token for authentication (can be empty)
	AuthToken string

	// JSON string to send as request body
	JSONBody string
}

// Response contains the response data from an HTTP request
type Response struct {
	// Response body as a byte slice
	Body []byte

	// HTTP status code
	StatusCode int
}

// PostJSON performs an HTTP POST request with JSON payload
//
// Makes an HTTP POST request to the specified URL with the given JSON body.
// Automatically sets Content-Type to application/json and includes
// Bearer authentication if AuthToken is provided.
func PostJSON(req *Request) (*Response, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create request
	httpReq, err := http.NewRequest("POST", req.URL, bytes.NewBufferString(req.JSONBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	
	// Add authorization header if token is provided
	if req.AuthToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+req.AuthToken)
	}

	// Perform request
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Create response
	response := &Response{
		Body:       body,
		StatusCode: resp.StatusCode,
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		return response, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(body))
	}

	return response, nil
}