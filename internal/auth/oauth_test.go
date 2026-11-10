package auth

import (
	"testing"
)

func TestGeneratePKCE(t *testing.T) {
	verifier, challenge, err := GeneratePKCE()
	if err != nil {
		t.Fatalf("GeneratePKCE() error = %v", err)
	}

	if verifier == "" {
		t.Error("verifier should not be empty")
	}

	if challenge == "" {
		t.Error("challenge should not be empty")
	}

	// Verifier should be 43 characters (32 bytes base64url encoded)
	if len(verifier) != 43 {
		t.Errorf("verifier length = %d, want 43", len(verifier))
	}

	// Challenge should be 43 characters (32 bytes SHA256 base64url encoded)
	if len(challenge) != 43 {
		t.Errorf("challenge length = %d, want 43", len(challenge))
	}

	// Verify uniqueness
	verifier2, challenge2, err := GeneratePKCE()
	if err != nil {
		t.Fatalf("GeneratePKCE() second call error = %v", err)
	}

	if verifier == verifier2 {
		t.Error("verifier should be unique")
	}
	if challenge == challenge2 {
		t.Error("challenge should be unique")
	}
}

func TestGenerateState(t *testing.T) {
	state1, err := GenerateState()
	if err != nil {
		t.Fatalf("GenerateState() error = %v", err)
	}

	if state1 == "" {
		t.Error("state should not be empty")
	}

	// State should be 32 hex characters (16 bytes)
	if len(state1) != 32 {
		t.Errorf("state length = %d, want 32", len(state1))
	}

	// Verify uniqueness
	state2, err := GenerateState()
	if err != nil {
		t.Fatalf("GenerateState() second call error = %v", err)
	}

	if state1 == state2 {
		t.Error("state should be unique")
	}
}

func TestValidateState(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		actual   string
		want     bool
	}{
		{"matching states", "abc123", "abc123", true},
		{"different states", "abc123", "xyz789", false},
		{"different lengths", "abc", "abcd", false},
		{"empty expected", "", "abc", false},
		{"empty actual", "abc", "", false},
		{"both empty", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateState(tt.expected, tt.actual)
			if got != tt.want {
				t.Errorf("ValidateState(%q, %q) = %v, want %v", tt.expected, tt.actual, got, tt.want)
			}
		})
	}
}
