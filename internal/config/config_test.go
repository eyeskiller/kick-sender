package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigSaveAndLoad(t *testing.T) {
	// Create temp directory for test
	tmpDir, err := os.MkdirTemp("", "kick-interval-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Override home directory for testing
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cfg := &Config{
		Channel:           "testchannel",
		BroadcasterUserID: 12345,
		Message:           "Hello world",
		Interval:          "10m",
		SendImmediately:   false,
	}

	// Save
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify file exists
	configPath := filepath.Join(tmpDir, "Library", "Application Support", "kick-interval", "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("config file was not created")
	}

	// Load
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.Channel != cfg.Channel {
		t.Errorf("Channel = %q, want %q", loaded.Channel, cfg.Channel)
	}
	if loaded.BroadcasterUserID != cfg.BroadcasterUserID {
		t.Errorf("BroadcasterUserID = %d, want %d", loaded.BroadcasterUserID, cfg.BroadcasterUserID)
	}
	if loaded.Message != cfg.Message {
		t.Errorf("Message = %q, want %q", loaded.Message, cfg.Message)
	}
	if loaded.Interval != cfg.Interval {
		t.Errorf("Interval = %q, want %q", loaded.Interval, cfg.Interval)
	}
}

func TestHasSettings(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want bool
	}{
		{
			"complete settings",
			Config{Channel: "test", Message: "hello", Interval: "10m"},
			true,
		},
		{
			"missing channel",
			Config{Message: "hello", Interval: "10m"},
			false,
		},
		{
			"missing message",
			Config{Channel: "test", Interval: "10m"},
			false,
		},
		{
			"missing interval",
			Config{Channel: "test", Message: "hello"},
			false,
		},
		{
			"empty config",
			Config{},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.HasSettings()
			if got != tt.want {
				t.Errorf("HasSettings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadNonExistent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kick-interval-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Channel != "" || cfg.Message != "" || cfg.Interval != "" {
		t.Error("expected empty config for non-existent file")
	}
}
