package kick

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

// NormalizeChannel normalizes a channel input to a slug
// Accepts either a slug or a full URL like https://kick.com/username
func NormalizeChannel(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("channel cannot be empty")
	}

	// If it's a URL, extract the path
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		parsed, err := url.Parse(input)
		if err != nil {
			return "", fmt.Errorf("invalid URL: %w", err)
		}

		// Extract slug from path
		path := strings.Trim(parsed.Path, "/")
		parts := strings.Split(path, "/")
		if len(parts) == 0 || parts[0] == "" {
			return "", fmt.Errorf("invalid Kick URL: no channel found")
		}

		return strings.ToLower(parts[0]), nil
	}

	// It's already a slug
	return strings.ToLower(input), nil
}

// ValidateSlug validates a channel slug
func ValidateSlug(slug string) error {
	if slug == "" {
		return fmt.Errorf("slug cannot be empty")
	}

	if len(slug) > 25 {
		return fmt.Errorf("slug too long (max 25 characters)")
	}

	// Check for valid characters (alphanumeric, hyphens, underscores)
	for _, r := range slug {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_') {
			return fmt.Errorf("invalid character in slug: %c", r)
		}
	}

	return nil
}

// ValidateMessage validates a chat message according to Kick's requirements
// - Maximum 500 Unicode grapheme clusters
// - Maximum 2048 bytes when encoded in UTF-8
func ValidateMessage(content string) error {
	if content == "" {
		return fmt.Errorf("message cannot be empty")
	}

	// Count grapheme clusters (simplified: count runes as a reasonable approximation)
	// A proper implementation would use golang.org/x/text/unicode/norm
	runeCount := utf8.RuneCountInString(content)
	if runeCount > 500 {
		return fmt.Errorf("message too long (max 500 characters, got %d)", runeCount)
	}

	// Check byte length
	byteLen := len([]byte(content))
	if byteLen > 2048 {
		return fmt.Errorf("message too large (max 2048 bytes, got %d)", byteLen)
	}

	return nil
}

// GraphemeCount returns an approximate count of grapheme clusters
// This is a simplified implementation using rune count
func GraphemeCount(s string) int {
	return utf8.RuneCountInString(s)
}
