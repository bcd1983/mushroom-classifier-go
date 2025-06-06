// Package openai provides OpenAI API integration for image analysis
package openai

import (
	"encoding/json"
	"fmt"

	"github.com/mushroom-classifier/mushroom-classifier-go/httpclient"
)

// Request contains parameters for an OpenAI API image analysis request
type Request struct {
	// API key for authentication
	APIKey string

	// Full URL to the OpenAI API endpoint
	APIURL string

	// Model identifier (e.g., "gpt-4o")
	Model string

	// Text prompt describing what to analyze
	Prompt string

	// Base64 encoded image data (optional)
	Base64Image string

	// Maximum tokens in the response
	MaxTokens int
}

// Response contains the result from OpenAI API call
type Response struct {
	// Response content from the model (valid if Success=true)
	Content string

	// Error message (valid if Success=false)
	ErrorMessage string

	// Success flag: true for success, false for failure
	Success bool
}

// chatCompletionRequest represents the JSON structure for OpenAI API request
type chatCompletionRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
	MaxTokens int      `json:"max_tokens"`
}

// message represents a chat message in the OpenAI API
type message struct {
	Role    string    `json:"role"`
	Content []content `json:"content"`
}

// content represents the content of a message (text or image)
type content struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *imageURL `json:"image_url,omitempty"`
}

// imageURL represents an image URL in the OpenAI API
type imageURL struct {
	URL string `json:"url"`
}

// chatCompletionResponse represents the JSON structure for OpenAI API response
type chatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

// AnalyzeImage sends an image along with a text prompt to OpenAI's API for analysis
//
// The function handles all API communication, request formatting, and
// response parsing. If Base64Image is empty, only the text prompt is sent.
func AnalyzeImage(req *Request) (*Response, error) {
	// Validate request
	if req.APIKey == "" {
		return &Response{
			Success:      false,
			ErrorMessage: "API key is required",
		}, nil
	}

	if req.APIURL == "" {
		return &Response{
			Success:      false,
			ErrorMessage: "API URL is required",
		}, nil
	}

	if req.Prompt == "" {
		return &Response{
			Success:      false,
			ErrorMessage: "Prompt is required",
		}, nil
	}

	// Set defaults
	if req.Model == "" {
		req.Model = "gpt-4o"
	}

	if req.MaxTokens <= 0 {
		req.MaxTokens = 1000
	}

	// Build message content
	messageContent := []content{
		{
			Type: "text",
			Text: req.Prompt,
		},
	}

	// Add image if provided
	if req.Base64Image != "" {
		messageContent = append(messageContent, content{
			Type: "image_url",
			ImageURL: &imageURL{
				URL: fmt.Sprintf("data:image/jpeg;base64,%s", req.Base64Image),
			},
		})
	}

	// Build request
	chatReq := chatCompletionRequest{
		Model: req.Model,
		Messages: []message{
			{
				Role:    "user",
				Content: messageContent,
			},
		},
		MaxTokens: req.MaxTokens,
	}

	// Marshal to JSON
	jsonBody, err := json.Marshal(chatReq)
	if err != nil {
		return &Response{
			Success:      false,
			ErrorMessage: fmt.Sprintf("Failed to marshal request: %v", err),
		}, nil
	}

	// Make HTTP request
	httpReq := &httpclient.Request{
		URL:       req.APIURL,
		AuthToken: req.APIKey,
		JSONBody:  string(jsonBody),
	}

	httpResp, err := httpclient.PostJSON(httpReq)
	if err != nil {
		return &Response{
			Success:      false,
			ErrorMessage: fmt.Sprintf("HTTP request failed: %v", err),
		}, nil
	}

	// Parse response
	var chatResp chatCompletionResponse
	if err := json.Unmarshal(httpResp.Body, &chatResp); err != nil {
		return &Response{
			Success:      false,
			ErrorMessage: fmt.Sprintf("Failed to parse response: %v", err),
		}, nil
	}

	// Check for API error
	if chatResp.Error != nil {
		return &Response{
			Success:      false,
			ErrorMessage: fmt.Sprintf("OpenAI API error: %s", chatResp.Error.Message),
		}, nil
	}

	// Extract content from response
	if len(chatResp.Choices) == 0 {
		return &Response{
			Success:      false,
			ErrorMessage: "No response from OpenAI API",
		}, nil
	}

	return &Response{
		Success: true,
		Content: chatResp.Choices[0].Message.Content,
	}, nil
}