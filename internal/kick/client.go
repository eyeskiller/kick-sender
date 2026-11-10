package kick

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	BaseURL = "https://api.kick.com"
	AuthURL = "https://id.kick.com"
)

// Client is the Kick API client
type Client struct {
	httpClient   *http.Client
	tokenGetter  func() string
	tokenRefresher func() error
}

// NewClient creates a new Kick API client
func NewClient(tokenGetter func() string, tokenRefresher func() error) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		tokenGetter: tokenGetter,
		tokenRefresher: tokenRefresher,
	}
}

// Channel represents a Kick channel
type Channel struct {
	Slug              string `json:"slug"`
	BroadcasterUserID int64  `json:"broadcaster_user_id"`
	StreamTitle       string `json:"stream_title,omitempty"`
	IsLive            bool   `json:"is_live,omitempty"`
}

// GetChannel retrieves channel information by slug
func (c *Client) GetChannel(slug string) (*Channel, error) {
	return c.getChannel(slug, false)
}

func (c *Client) getChannel(slug string, retried bool) (*Channel, error) {
	reqURL := fmt.Sprintf("%s/public/v1/channels?slug=%s", BaseURL, url.QueryEscape(slug))

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.tokenGetter())
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		if retried {
			return nil, fmt.Errorf("unauthorized after token refresh - re-authentication required")
		}
		if c.tokenRefresher != nil {
			if err := c.tokenRefresher(); err == nil {
				return c.getChannel(slug, true)
			} else {
				return nil, fmt.Errorf("token refresh failed: %w", err)
			}
		}
		return nil, fmt.Errorf("unauthorized: %s", string(body))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Data []Channel `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("channel not found: %s", slug)
	}

	return &response.Data[0], nil
}

// ChatMessage represents a chat message to send
type ChatMessage struct {
	Type              string `json:"type"`
	BroadcasterUserID int64  `json:"broadcaster_user_id"`
	Content           string `json:"content"`
}

// ChatResponse represents the response from sending a chat message
type ChatResponse struct {
	IsSent    bool   `json:"is_sent"`
	MessageID string `json:"message_id,omitempty"`
}

// SendChatMessage sends a chat message to a channel
func (c *Client) SendChatMessage(broadcasterUserID int64, content string) (*ChatResponse, error) {
	return c.sendChatMessage(broadcasterUserID, content, false)
}

func (c *Client) sendChatMessage(broadcasterUserID int64, content string, retried bool) (*ChatResponse, error) {
	msg := ChatMessage{
		Type:              "user",
		BroadcasterUserID: broadcasterUserID,
		Content:           content,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequest("POST", BaseURL+"/public/v1/chat", bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.tokenGetter())
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		if retried {
			return nil, fmt.Errorf("unauthorized after token refresh - re-authentication required")
		}
		if c.tokenRefresher != nil {
			if err := c.tokenRefresher(); err == nil {
				return c.sendChatMessage(broadcasterUserID, content, true)
			} else {
				return nil, fmt.Errorf("token refresh failed: %w", err)
			}
		}
		return nil, fmt.Errorf("unauthorized: %s", string(body))
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("rate limited")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Data    ChatResponse `json:"data"`
		Message string       `json:"message"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !response.Data.IsSent {
		return nil, fmt.Errorf("message not sent: %s", response.Message)
	}

	return &response.Data, nil
}
