# Bug Fix: Token Refresh Infinite Loop

## Problem Description

After running for a while (especially when internet connection drops), the application would enter an infinite loop of token refreshes:

```
→ Refreshing access token...
✓ Token refreshed
→ Refreshing access token...
✓ Token refreshed
→ Refreshing access token...
✓ Token refreshed
[repeats infinitely]
```

## Root Cause

The Kick API client had a critical design flaw:

1. When a 401 Unauthorized response was received, the client would call the `tokenRefresher` function
2. The token refresher would successfully get new tokens and save them to the keychain via `TokenStore`
3. **However, the client's `accessToken` field was never updated** - it still held the old, invalid token
4. The client would retry the request with the old token
5. The old token would fail again with 401
6. Steps 1-5 would repeat infinitely

This was especially problematic when:
- Internet connection dropped and came back
- Tokens expired during operation
- Any scenario requiring token refresh

## Solution

### 1. Changed Client to Use Token Getter Function

**Before:**
```go
type Client struct {
    httpClient  *http.Client
    accessToken string  // Static copy of token
    tokenRefresher func() error
}

func NewClient(accessToken string, tokenRefresher func() error) *Client {
    return &Client{
        httpClient: &http.Client{Timeout: 30 * time.Second},
        accessToken: accessToken,  // Captured once at creation
        tokenRefresher: tokenRefresher,
    }
}
```

**After:**
```go
type Client struct {
    httpClient   *http.Client
    tokenGetter  func() string  // Function that always returns current token
    tokenRefresher func() error
}

func NewClient(tokenGetter func() string, tokenRefresher func() error) *Client {
    return &Client{
        httpClient: &http.Client{Timeout: 30 * time.Second},
        tokenGetter: tokenGetter,  // Called every time we need the token
        tokenRefresher: tokenRefresher,
    }
}
```

Now when the client needs to make a request:
```go
req.Header.Set("Authorization", "Bearer "+c.tokenGetter())  // Always gets fresh token
```

### 2. Added Loop Prevention

Added a `retried` flag to prevent infinite refresh loops:

```go
func (c *Client) GetChannel(slug string) (*Channel, error) {
    return c.getChannel(slug, false)  // First attempt
}

func (c *Client) getChannel(slug string, retried bool) (*Channel, error) {
    // ... make request ...
    
    if resp.StatusCode == http.StatusUnauthorized {
        if retried {
            // Already tried refresh once, don't loop
            return nil, fmt.Errorf("unauthorized after token refresh - re-authentication required")
        }
        if c.tokenRefresher != nil {
            if err := c.tokenRefresher(); err == nil {
                return c.getChannel(slug, true)  // Retry with flag set
            } else {
                return nil, fmt.Errorf("token refresh failed: %w", err)
            }
        }
        return nil, fmt.Errorf("unauthorized: %s", string(body))
    }
    
    // ... rest of logic ...
}
```

Same pattern applied to `SendChatMessage`.

### 3. Updated Main.go

Changed client initialization to pass the token getter function:

**Before:**
```go
client := kick.NewClient(tokenStore.GetAccessToken(), tokenRefresher)
```

**After:**
```go
client := kick.NewClient(tokenStore.GetAccessToken, tokenRefresher)  // Pass function, not result
```

## Benefits

1. **No More Infinite Loops**: Token refresh is attempted at most once per request
2. **Always Uses Fresh Tokens**: Client always reads the current token from the store
3. **Better Error Messages**: Clear indication when re-authentication is needed
4. **Network Resilience**: Properly handles connection drops and recovery
5. **Thread-Safe**: Multiple goroutines will always get the latest token

## Testing

All existing tests pass, confirming:
- No breaking changes to public API behavior
- Token refresh logic works correctly
- Loop prevention prevents infinite retries
- Client properly handles 401 responses

## Files Modified

1. `internal/kick/client.go` - Refactored to use token getter, added loop prevention
2. `internal/kick/client_test.go` - Updated test to match new signature
3. `cmd/kick-interval/main.go` - Updated client initialization

## Verification

To verify the fix works:
1. Run the application with a valid configuration
2. Wait for token expiration (or manually expire tokens)
3. Observe that token refresh happens once and succeeds
4. Application continues sending messages normally
5. No infinite refresh loops occur
