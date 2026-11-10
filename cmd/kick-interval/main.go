package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/user/kick-interval/internal/auth"
	"github.com/user/kick-interval/internal/config"
	"github.com/user/kick-interval/internal/kick"
	"github.com/user/kick-interval/internal/scheduler"
	"github.com/user/kick-interval/internal/server"
	"github.com/user/kick-interval/internal/ui"
)

var (
	version = "0.1.0"
)

func main() {
	// Define flags
	channelFlag := flag.String("channel", "", "Kick channel (slug or URL)")
	messageFlag := flag.String("message", "", "Message to send")
	intervalFlag := flag.String("interval", "", "Send interval (e.g., 60s, 5m, 10m, 1h)")
	sendNowFlag := flag.Bool("send-now", false, "Send first message immediately")
	dryRunFlag := flag.Bool("dry-run", false, "Validate without sending")
	versionFlag := flag.Bool("version", false, "Show version")
	helpFlag := flag.Bool("help", false, "Show help")

	flag.Usage = printUsage
	flag.Parse()

	// Handle subcommands
	args := flag.Args()
	if len(args) > 0 {
		switch args[0] {
		case "logout":
			handleLogout()
			return
		case "version":
			fmt.Printf("kick-interval %s\n", version)
			return
		case "serve":
			port := "8765"
			if len(args) > 1 {
				port = args[1]
			}
			if err := runServer(port); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			return
		}
	}

	if *helpFlag {
		printUsage()
		return
	}

	if *versionFlag {
		fmt.Printf("kick-interval %s\n", version)
		return
	}

	// Run main application
	if err := run(*channelFlag, *messageFlag, *intervalFlag, *sendNowFlag, *dryRunFlag); err != nil {
		ui.PrintError(err.Error())
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf("kick-interval v%s\n\n", version)
	fmt.Println("Automated Kick.com chat message sender")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  kick-interval                          Interactive mode")
	fmt.Println("  kick-interval --channel <slug> --message <msg> --interval <interval>")
	fmt.Println("  kick-interval --send-now               Send first message immediately")
	fmt.Println("  kick-interval --dry-run                Validate without sending")
	fmt.Println("  kick-interval logout                   Clear stored credentials")
	fmt.Println("  kick-interval serve [port]             Start HTTP server for GUI")
	fmt.Println("  kick-interval --version                Show version")
	fmt.Println("  kick-interval --help                   Show this help")
	fmt.Println()
	fmt.Println("Interval format: 60s, 5m, 10m, 1h, 2h")
}

func run(channelOverride, messageOverride, intervalOverride string, sendNow, dryRun bool) error {
	ui.PrintBanner()

	// Load saved configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Apply CLI overrides
	if channelOverride != "" {
		cfg.Channel = channelOverride
	}
	if messageOverride != "" {
		cfg.Message = messageOverride
	}
	if intervalOverride != "" {
		cfg.Interval = intervalOverride
	}
	if sendNow {
		cfg.SendImmediately = true
	}

	cliOverrides := channelOverride != "" || messageOverride != "" || intervalOverride != ""

	// Initialize token store
	tokenStore := &auth.TokenStore{}
	hasTokens := tokenStore.HasTokens()

	if !hasTokens {
		ui.PrintInfo("No authentication found")
		if err := handleAuthSetup(tokenStore); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
	}

	// Try to refresh token if needed
	if err := tokenStore.LoadTokens(); err != nil {
		return fmt.Errorf("failed to load tokens: %w", err)
	}

	// Create OAuth manager for token refresh
	oauthConfig, err := loadOAuthConfig()
	if err != nil {
		return fmt.Errorf("failed to load OAuth config: %w", err)
	}

	oauthMgr := auth.NewOAuthManager(oauthConfig, tokenStore)
	defer oauthMgr.Stop()

	tokenRefresher := func() error {
		ui.PrintInfo("Refreshing access token...")
		if err := oauthMgr.RefreshToken(); err != nil {
			return err
		}
		ui.PrintSuccess("Token refreshed")
		return nil
	}

	// Create Kick client
	client := kick.NewClient(tokenStore.GetAccessToken, tokenRefresher)

	// Check if we have saved settings and no overrides
	if cfg.HasSettings() && !cliOverrides {
		choice := handleSavedConfig(cfg)
		switch choice {
		case 0: // Start with saved settings
			// Continue with saved config
		case 1: // Change channel
			channel, err := promptChannel(client)
			if err != nil {
				return err
			}
			cfg.Channel = channel
			cfg.Save()
		case 2: // Change message
			cfg.Message = ui.Prompt("Enter the message to send:")
			if err := kick.ValidateMessage(cfg.Message); err != nil {
				return err
			}
			cfg.Save()
		case 3: // Change interval
			cfg.Interval = ui.Prompt("Choose sending interval (e.g., 10m):")
			if _, err := scheduler.ParseInterval(cfg.Interval); err != nil {
				return err
			}
			cfg.Save()
		case 4: // Change account
			if err := handleAuthSetup(tokenStore); err != nil {
				return err
			}
		case 5: // Exit
			fmt.Println("Goodbye!")
			return nil
		default:
			return fmt.Errorf("invalid choice")
		}
	}

	// Interactive mode for missing fields
	if cfg.Channel == "" {
		channel, err := promptChannel(client)
		if err != nil {
			return err
		}
		cfg.Channel = channel
	}

	// Validate channel
	ui.PrintInfo("Validating channel...")
	channelInfo, err := client.GetChannel(cfg.Channel)
	if err != nil {
		return fmt.Errorf("channel validation failed: %w", err)
	}
	cfg.BroadcasterUserID = channelInfo.BroadcasterUserID
	ui.PrintSuccess(fmt.Sprintf("Channel found: %s (ID: %d)", channelInfo.Slug, channelInfo.BroadcasterUserID))
	if channelInfo.IsLive {
		ui.PrintInfo("Channel is LIVE")
	}

	if cfg.Message == "" {
		cfg.Message = ui.Prompt("Enter the message to send:")
		if err := kick.ValidateMessage(cfg.Message); err != nil {
			return err
		}
	}

	if cfg.Interval == "" {
		cfg.Interval = ui.Prompt("Choose sending interval (e.g., 60s, 10m, 1h):")
	}

	// Parse interval
	interval, err := scheduler.ParseInterval(cfg.Interval)
	if err != nil {
		return fmt.Errorf("invalid interval: %w", err)
	}

	if interval < 60*time.Second {
		return fmt.Errorf("interval must be at least 60 seconds")
	}

	// Get username for display
	username := getAuthenticatedUsername(client)

	// Show configuration
	ui.PrintConfig(username, cfg.Channel, cfg.Message, scheduler.FormatInterval(interval))

	// Save configuration
	if err := cfg.Save(); err != nil {
		ui.PrintWarning(fmt.Sprintf("Could not save config: %v", err))
	}

	if dryRun {
		ui.PrintSuccess("Dry run complete - configuration is valid")
		return nil
	}

	// Confirm
	if !ui.PromptYesNo("Start sending?") {
		fmt.Println("Cancelled.")
		return nil
	}

	// Start scheduler
	return startScheduler(client, cfg, interval, sendNow)
}

func handleAuthSetup(tokenStore *auth.TokenStore) error {
	fmt.Println()
	fmt.Println("To use this application, you need a Kick Developer App.")
	fmt.Println("Create one at: https://kick.com/settings/developer")
	fmt.Println()
	fmt.Println("You'll need:")
	fmt.Println("  - Client ID")
	fmt.Println("  - Client Secret")
	fmt.Println("  - Redirect URL (set to: http://localhost:8765/callback)")
	fmt.Println()

	clientID := ui.Prompt("Enter your Client ID:")
	if clientID == "" {
		return fmt.Errorf("client ID required")
	}

	clientSecret := ui.Prompt("Enter your Client Secret:")
	if clientSecret == "" {
		return fmt.Errorf("client secret required")
	}

	redirectURI := ui.Prompt("Enter redirect URI (default: http://localhost:8765/callback):")
	if redirectURI == "" {
		redirectURI = "http://localhost:8765/callback"
	}

	// Save OAuth config
	oauthConfig := auth.OAuthConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
	}
	if err := saveOAuthConfig(oauthConfig); err != nil {
		return fmt.Errorf("failed to save OAuth config: %w", err)
	}

	// Start OAuth flow
	oauthMgr := auth.NewOAuthManager(oauthConfig, tokenStore)
	defer oauthMgr.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := oauthMgr.StartAuthFlow(ctx); err != nil {
		return err
	}

	ui.PrintSuccess("Authentication successful!")
	return nil
}

func promptChannel(client *kick.Client) (string, error) {
	for {
		input := ui.Prompt("Enter Kick channel (slug or URL):")
		slug, err := kick.NormalizeChannel(input)
		if err != nil {
			ui.PrintError(err.Error())
			continue
		}

		if err := kick.ValidateSlug(slug); err != nil {
			ui.PrintError(err.Error())
			continue
		}

		// Validate channel exists
		ui.PrintInfo("Looking up channel...")
		_, err = client.GetChannel(slug)
		if err != nil {
			ui.PrintError(fmt.Sprintf("Channel not found: %v", err))
			continue
		}

		return slug, nil
	}
}

func handleSavedConfig(cfg *config.Config) int {
	fmt.Println("Saved configuration found:")
	fmt.Println()
	fmt.Printf("Channel:  %s\n", cfg.Channel)
	fmt.Printf("Message:  %s\n", cfg.Message)
	fmt.Printf("Interval: %s\n", cfg.Interval)
	fmt.Println()

	choices := []string{
		"Start with saved settings",
		"Change channel",
		"Change message",
		"Change interval",
		"Change Kick account",
		"Exit",
	}

	return ui.PromptChoice("What would you like to do?", choices)
}

func getAuthenticatedUsername(client *kick.Client) string {
	// Get own channel (no parameters returns authenticated user's channel)
	// Since we don't have a dedicated user info endpoint, use a placeholder
	// In production, this could use the user:read scope endpoint
	return "authenticated user"
}

func startScheduler(client *kick.Client, cfg *config.Config, interval time.Duration, sendNow bool) error {
	var sentCount int64
	var failedCount int64

	sendFunc := func() error {
		_, err := client.SendChatMessage(cfg.BroadcasterUserID, cfg.Message)
		return err
	}

	onSend := func(t time.Time, err error) {
		if err != nil {
			atomic.AddInt64(&failedCount, 1)
		} else {
			atomic.AddInt64(&sentCount, 1)
		}
		ui.PrintSendResult(t, err)
		ui.PrintSummary(int(atomic.LoadInt64(&sentCount)), int(atomic.LoadInt64(&failedCount)))
	}

	sched := scheduler.NewScheduler(interval, sendFunc, onSend, sendNow)

	ui.PrintRunning("authenticated user", cfg.Channel, scheduler.FormatInterval(interval))
	nextSend := sched.NextSend()
	if sendNow {
		ui.PrintInfo(fmt.Sprintf("Sending first message now..."))
	} else {
		ui.PrintInfo(fmt.Sprintf("Next message: %s", nextSend.Format("15:04:05")))
	}

	sched.Start()

	// Wait for interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("\nPress Ctrl+C to stop")

	<-sigChan

	fmt.Println("\n\nShutting down...")
	sched.Stop()

	ui.PrintSummary(int(atomic.LoadInt64(&sentCount)), int(atomic.LoadInt64(&failedCount)))
	ui.PrintSuccess("Stopped")

	return nil
}

func handleLogout() {
	tokenStore := &auth.TokenStore{}
	if err := tokenStore.ClearTokens(); err != nil {
		ui.PrintError(fmt.Sprintf("Failed to clear credentials: %v", err))
		os.Exit(1)
	}
	ui.PrintSuccess("Credentials cleared")
}

func loadOAuthConfig() (auth.OAuthConfig, error) {
	configDir, err := os.UserHomeDir()
	if err != nil {
		return auth.OAuthConfig{}, err
	}
	configDir = configDir + "/Library/Application Support/kick-interval"

	data, err := os.ReadFile(configDir + "/oauth.json")
	if err != nil {
		return auth.OAuthConfig{}, fmt.Errorf("no OAuth config found. Run the application to set up authentication")
	}

	var cfg auth.OAuthConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return auth.OAuthConfig{}, err
	}

	return cfg, nil
}

func saveOAuthConfig(cfg auth.OAuthConfig) error {
	configDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configDir = configDir + "/Library/Application Support/kick-interval"

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configDir+"/oauth.json", data, 0600)
}

func runServer(port string) error {
	// Load saved configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize token store
	tokenStore := &auth.TokenStore{}
	// Load existing tokens from keychain (ignore error if none exist yet)
	_ = tokenStore.LoadTokens()

	// Load OAuth config
	oauthConfig, err := loadOAuthConfig()
	if err != nil {
		// OAuth config not found, create default
		oauthConfig = auth.OAuthConfig{
			RedirectURI: "http://localhost:8765/callback",
		}
	}

	// Create and start server
	srv := server.NewServer(cfg, tokenStore, oauthConfig)

	fmt.Printf("Kick Interval server starting on port %s...\n", port)
	fmt.Println("API available at http://localhost:" + port)
	fmt.Println("Press Ctrl+C to stop")

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.Start(port)
	}()

	select {
	case err := <-errChan:
		return err
	case <-sigChan:
		fmt.Println("\nShutting down...")
		return srv.Stop()
	}
}
