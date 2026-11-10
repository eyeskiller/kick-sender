package auth

import (
	"fmt"
	"os/exec"

	"github.com/zalando/go-keyring"
)

const (
	serviceName   = "kick-interval"
	accessTokenKey = "access_token"
	refreshTokenKey = "refresh_token"
)

// TokenStore manages OAuth tokens in the macOS Keychain
type TokenStore struct {
	accessToken  string
	refreshToken string
}

// StoreTokens saves tokens to the keychain
func (ts *TokenStore) StoreTokens(accessToken, refreshToken string) error {
	if err := keyring.Set(serviceName, accessTokenKey, accessToken); err != nil {
		return fmt.Errorf("failed to store access token: %w", err)
	}
	if err := keyring.Set(serviceName, refreshTokenKey, refreshToken); err != nil {
		return fmt.Errorf("failed to store refresh token: %w", err)
	}
	ts.accessToken = accessToken
	ts.refreshToken = refreshToken
	return nil
}

// LoadTokens retrieves tokens from the keychain
func (ts *TokenStore) LoadTokens() error {
	accessToken, err := keyring.Get(serviceName, accessTokenKey)
	if err != nil {
		if err == keyring.ErrNotFound {
			return fmt.Errorf("no tokens found")
		}
		return fmt.Errorf("failed to load access token: %w", err)
	}

	refreshToken, err := keyring.Get(serviceName, refreshTokenKey)
	if err != nil {
		if err == keyring.ErrNotFound {
			return fmt.Errorf("no refresh token found")
		}
		return fmt.Errorf("failed to load refresh token: %w", err)
	}

	ts.accessToken = accessToken
	ts.refreshToken = refreshToken
	return nil
}

// GetAccessToken returns the current access token
func (ts *TokenStore) GetAccessToken() string {
	return ts.accessToken
}

// GetRefreshToken returns the current refresh token
func (ts *TokenStore) GetRefreshToken() string {
	return ts.refreshToken
}

// ClearTokens removes all tokens from the keychain
func (ts *TokenStore) ClearTokens() error {
	if err := keyring.Delete(serviceName, accessTokenKey); err != nil && err != keyring.ErrNotFound {
		return fmt.Errorf("failed to delete access token: %w", err)
	}
	if err := keyring.Delete(serviceName, refreshTokenKey); err != nil && err != keyring.ErrNotFound {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}
	ts.accessToken = ""
	ts.refreshToken = ""
	return nil
}

// HasTokens checks if tokens exist in the keychain
func (ts *TokenStore) HasTokens() bool {
	_, err := keyring.Get(serviceName, accessTokenKey)
	return err == nil
}

// OpenBrowser opens the default browser to the specified URL
func OpenBrowser(url string) error {
	return exec.Command("open", url).Start()
}
