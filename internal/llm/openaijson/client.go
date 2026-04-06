package openaijson

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// DefaultBaseURL is the OpenAI API v1 base URL.
const DefaultBaseURL = "https://api.openai.com/v1"

// Client performs chat.completions requests with response_format json_object.
type Client struct {
	HTTP    *http.Client
	BaseURL string
}

func (c *Client) base() string {
	if strings.TrimSpace(c.BaseURL) == "" {
		return DefaultBaseURL
	}
	return strings.TrimRight(strings.TrimSpace(c.BaseURL), "/")
}

// ChatCompletionJSON returns the assistant message content (expected JSON object text).
func (c *Client) ChatCompletionJSON(ctx context.Context, apiKey, model, systemPrompt, userContent string) (string, error) {
	if strings.TrimSpace(apiKey) == "" {
		return "", fmt.Errorf("openai: empty api key")
	}
	if strings.TrimSpace(model) == "" {
		return "", fmt.Errorf("openai: empty model")
	}

	httpClient := c.HTTP
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 90 * time.Second}
	}

	body := map[string]interface{}{
		"model": model,
		"response_format": map[string]string{
			"type": "json_object",
		},
		"temperature": 0.1,
		"max_tokens":  1024,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userContent},
		},
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("openai: marshal request: %w", err)
	}

	url := c.base() + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("openai: request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return "", fmt.Errorf("openai: read body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("openai: status %d: %s", resp.StatusCode, truncateForErr(string(respBody), 512))
	}

	var out chatCompletionsResponse
	if err := json.Unmarshal(respBody, &out); err != nil {
		return "", fmt.Errorf("openai: decode response: %w", err)
	}
	if len(out.Choices) == 0 || out.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("openai: empty choices in response")
	}
	return out.Choices[0].Message.Content, nil
}

type chatCompletionsResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func truncateForErr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
