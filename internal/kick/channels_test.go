package kick

import (
	"testing"
)

func TestNormalizeChannel(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"plain slug", "examplechannel", "examplechannel", false},
		{"uppercase slug", "ExampleChannel", "examplechannel", false},
		{"https URL", "https://kick.com/examplechannel", "examplechannel", false},
		{"http URL", "http://kick.com/examplechannel", "examplechannel", false},
		{"URL with trailing slash", "https://kick.com/examplechannel/", "examplechannel", false},
		{"URL with extra path", "https://kick.com/examplechannel/about", "examplechannel", false},
		{"slug with hyphen", "example-channel", "example-channel", false},
		{"slug with underscore", "example_channel", "example_channel", false},
		{"empty", "", "", true},
		{"whitespace only", "   ", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeChannel(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeChannel(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NormalizeChannel(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateSlug(t *testing.T) {
	tests := []struct {
		name    string
		slug    string
		wantErr bool
	}{
		{"valid simple", "examplechannel", false},
		{"valid with hyphen", "example-channel", false},
		{"valid with underscore", "example_channel", false},
		{"valid with numbers", "channel123", false},
		{"empty", "", true},
		{"too long", "abcdefghijklmnopqrstuvwxyz", true},
		{"invalid char space", "example channel", true},
		{"invalid char dot", "example.channel", true},
		{"invalid char at", "example@channel", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSlug(tt.slug)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSlug(%q) error = %v, wantErr %v", tt.slug, err, tt.wantErr)
			}
		})
	}
}

func TestValidateMessage(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{"valid simple", "Hello world", false},
		{"valid unicode", "Hello 世界", false},
		{"empty", "", true},
		{"too many characters", string(make([]rune, 501)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMessage(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMessage(%q) error = %v, wantErr %v", tt.content, err, tt.wantErr)
			}
		})
	}

	// Test byte limit
	t.Run("too many bytes", func(t *testing.T) {
		// Create a string with 2049 bytes
		content := make([]byte, 2049)
		for i := range content {
			content[i] = 'a'
		}
		err := ValidateMessage(string(content))
		if err == nil {
			t.Error("ValidateMessage should fail for messages > 2048 bytes")
		}
	})

	// Test exactly at limits
	t.Run("exactly 500 characters", func(t *testing.T) {
		content := make([]rune, 500)
		for i := range content {
			content[i] = 'a'
		}
		err := ValidateMessage(string(content))
		if err != nil {
			t.Errorf("ValidateMessage should pass for exactly 500 characters: %v", err)
		}
	})
}

func TestGraphemeCount(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"empty", "", 0},
		{"simple", "hello", 5},
		{"unicode", "hello 世界", 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GraphemeCount(tt.input)
			if got != tt.want {
				t.Errorf("GraphemeCount(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}
