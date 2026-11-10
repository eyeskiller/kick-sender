package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/user/kick-interval/internal/ui"
)

// OAuthConfig holds OAuth client configuration
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// OAuthManager handles the OAuth flow
type OAuthManager struct {
	config       OAuthConfig
	tokenStore   *TokenStore
	httpServer   *http.Server
	callbackCode string
	callbackState string
}

// NewOAuthManager creates a new OAuth manager
func NewOAuthManager(config OAuthConfig, tokenStore *TokenStore) *OAuthManager {
	return &OAuthManager{
		config:     config,
		tokenStore: tokenStore,
	}
}

// StartAuthFlow initiates the OAuth authorization flow
func (m *OAuthManager) StartAuthFlow(ctx context.Context) error {
	ui.PrintInfo("Starting authentication flow...")

	// Generate PKCE
	verifier, challenge, err := GeneratePKCE()
	if err != nil {
		return fmt.Errorf("failed to generate PKCE: %w", err)
	}

	// Generate state
	state, err := GenerateState()
	if err != nil {
		return fmt.Errorf("failed to generate state: %w", err)
	}

	// Build authorization URL
	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", m.config.ClientID)
	params.Set("redirect_uri", m.config.RedirectURI)
	params.Set("scope", "chat:write channel:read user:read")
	params.Set("code_challenge", challenge)
	params.Set("code_challenge_method", "S256")
	params.Set("state", state)

	authURL := fmt.Sprintf("https://id.kick.com/oauth/authorize?%s", params.Encode())

	// Start local HTTP server to receive callback
	m.callbackCode = ""
	m.callbackState = ""

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		returnedState := r.URL.Query().Get("state")

		if code == "" {
			fmt.Fprintf(w, "<html><body><h1>Error</h1><p>No authorization code received.</p></body></html>")
			return
		}

		if !ValidateState(state, returnedState) {
			fmt.Fprintf(w, "<html><body><h1>Error</h1><p>Invalid state parameter.</p></body></html>")
			return
		}

		m.callbackCode = code
		m.callbackState = returnedState

		fmt.Fprintf(w, "<html><body><h1>Success!</h1><p>You can close this window and return to the terminal.</p></body></html>")
	})

	// Parse redirect URI to get port
	parsedURL, err := url.Parse(m.config.RedirectURI)
	if err != nil {
		return fmt.Errorf("invalid redirect URI: %w", err)
	}

	addr := parsedURL.Host
	if !strings.Contains(addr, ":") {
		addr = addr + ":80"
	}

	m.httpServer = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		if err := m.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			ui.PrintError(fmt.Sprintf("HTTP server error: %v", err))
		}
	}()

	// Open browser
	ui.PrintInfo("Opening browser for authentication...")
	if err := OpenBrowser(authURL); err != nil {
		ui.PrintWarning("Could not open browser automatically")
		fmt.Printf("\nPlease open this URL manually:\n%s\n\n", authURL)
	}

	// Wait for callback
	ui.PrintInfo("Waiting for authorization...")
	timeout := time.After(5 * time.Minute)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout:
			return fmt.Errorf("authentication timeout")
		case <-ticker.C:
			if m.callbackCode != "" {
				// Exchange code for tokens
				return m.exchangeCodeForTokens(verifier)
			}
		}
	}
}

func (m *OAuthManager) exchangeCodeForTokens(verifier string) error {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", m.config.ClientID)
	data.Set("client_secret", m.config.ClientSecret)
	data.Set("redirect_uri", m.config.RedirectURI)
	data.Set("code_verifier", verifier)
	data.Set("code", m.callbackCode)

	req, err := http.NewRequest("POST", "https://id.kick.com/oauth/token", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send token request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token exchange failed (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return fmt.Errorf("no access token in response")
	}

	return m.tokenStore.StoreTokens(tokenResp.AccessToken, tokenResp.RefreshToken)
}

// RefreshToken refreshes the access token using the refresh token
func (m *OAuthManager) RefreshToken() error {
	refreshToken := m.tokenStore.GetRefreshToken()
	if refreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", m.config.ClientID)
	data.Set("client_secret", m.config.ClientSecret)
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", "https://id.kick.com/oauth/token", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create refresh request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send refresh request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read refresh response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token refresh failed (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse refresh response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return fmt.Errorf("no access token in refresh response")
	}

	return m.tokenStore.StoreTokens(tokenResp.AccessToken, tokenResp.RefreshToken)
}

// UpdateConfig updates the OAuth configuration (used by the onboarding flow)
func (m *OAuthManager) UpdateConfig(config OAuthConfig) {
	m.config = config
}

// GetConfig returns the current OAuth config
func (m *OAuthManager) GetConfig() OAuthConfig {
	return m.config
}

// HasConfig returns true if client ID and secret are set
func (m *OAuthManager) HasConfig() bool {
	return m.config.ClientID != "" && m.config.ClientSecret != ""
}

// Stop stops the OAuth HTTP server
func (m *OAuthManager) Stop() {
	if m.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		m.httpServer.Shutdown(ctx)
	}
}
