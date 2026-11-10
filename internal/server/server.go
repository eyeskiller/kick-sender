package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/user/kick-interval/internal/auth"
	"github.com/user/kick-interval/internal/config"
	"github.com/user/kick-interval/internal/kick"
	"github.com/user/kick-interval/internal/scheduler"
)

// Server is the HTTP API server
type Server struct {
	cfg         *config.Config
	tokenStore  *auth.TokenStore
	oauthMgr    *auth.OAuthManager
	client      *kick.Client
	scheduler   *scheduler.Scheduler
	schedMutex  sync.Mutex
	logs        []LogEntry
	logsMutex   sync.Mutex
	running     bool
	server      *http.Server
}

// LogEntry represents a log message
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := Status{
		Authenticated:  s.tokenStore.HasTokens(),
		Channel:        s.cfg.Channel,
		ChannelID:      s.cfg.BroadcasterUserID,
		Message:        s.cfg.Message,
		Interval:       s.cfg.Interval,
		Running:        s.running,
		HasOAuthConfig: s.oauthMgr.HasConfig(),
	}

	if status.Running && s.scheduler != nil {
		next := s.scheduler.NextSend()
		status.NextSend = next.Format("15:04:05")
		status.NextSendUnix = next.Unix()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// Status represents the current application status
type Status struct {
	Authenticated  bool   `json:"authenticated"`
	Username       string `json:"username"`
	Channel        string `json:"channel"`
	ChannelID      int64  `json:"channel_id"`
	Message        string `json:"message"`
	Interval       string `json:"interval"`
	Running        bool   `json:"running"`
	NextSend       string `json:"next_send,omitempty"`
	NextSendUnix   int64  `json:"next_send_unix,omitempty"`
	MessagesSent   int64  `json:"messages_sent"`
	FailedCount    int64  `json:"failed_count"`
	HasOAuthConfig bool   `json:"has_oauth_config"`
}

// NewServer creates a new HTTP server
func NewServer(cfg *config.Config, tokenStore *auth.TokenStore, oauthConfig auth.OAuthConfig) *Server {
	oauthMgr := auth.NewOAuthManager(oauthConfig, tokenStore)
	client := kick.NewClient(tokenStore.GetAccessToken, func() error {
		return oauthMgr.RefreshToken()
	})

	s := &Server{
		cfg:        cfg,
		tokenStore: tokenStore,
		oauthMgr:   oauthMgr,
		client:     client,
		logs:       make([]LogEntry, 0, 1000),
	}

	return s
}

// Start starts the HTTP server
func (s *Server) Start(port string) error {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/status", s.handleStatus)
	mux.HandleFunc("/api/auth/start", s.handleAuthStart)
	mux.HandleFunc("/api/auth/config", s.handleAuthConfig)
	mux.HandleFunc("/api/channel", s.handleChannel)
	mux.HandleFunc("/api/message", s.handleMessage)
	mux.HandleFunc("/api/interval", s.handleInterval)
	mux.HandleFunc("/api/start", s.handleStart)
	mux.HandleFunc("/api/stop", s.handleStop)
	mux.HandleFunc("/api/logs", s.handleLogs)
	mux.HandleFunc("/api/reset", s.handleReset)

	// CORS middleware for Electron
	handler := corsMiddleware(mux)

	s.server = &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	s.addLog("INFO", fmt.Sprintf("Server starting on port %s", port))

	return s.server.ListenAndServe()
}

// Stop stops the HTTP server
func (s *Server) Stop() error {
	if s.scheduler != nil {
		s.scheduler.Stop()
	}
	if s.oauthMgr != nil {
		s.oauthMgr.Stop()
	}
	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.server.Shutdown(ctx)
	}
	return nil
}

func (s *Server) addLog(level, message string) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
	}

	s.logsMutex.Lock()
	defer s.logsMutex.Unlock()

	s.logs = append(s.logs, entry)

	// Keep only last 1000 logs
	if len(s.logs) > 1000 {
		s.logs = s.logs[len(s.logs)-1000:]
	}
}

func (s *Server) handleAuthStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.tokenStore.HasTokens() {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "already_authenticated"})
		return
	}

	if !s.oauthMgr.HasConfig() {
		http.Error(w, "OAuth config not set. Provide client_id and client_secret first via /api/auth/config", http.StatusBadRequest)
		return
	}

	s.addLog("INFO", "Starting OAuth flow")

	// Start OAuth flow in background
	go func() {
		ctx := context.Background()
		if err := s.oauthMgr.StartAuthFlow(ctx); err != nil {
			s.addLog("ERROR", fmt.Sprintf("Authentication failed: %v", err))
		} else {
			s.addLog("INFO", "Authentication successful")
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "started",
		"message": "Authentication flow started. Browser should open automatically.",
	})
}

func (s *Server) handleAuthConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		RedirectURI  string `json:"redirect_uri"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ClientID == "" || req.ClientSecret == "" {
		http.Error(w, "client_id and client_secret are required", http.StatusBadRequest)
		return
	}

	if req.RedirectURI == "" {
		req.RedirectURI = "http://localhost:8765/callback"
	}

	// Persist to disk so it survives restarts
	if err := saveOAuthConfigFile(auth.OAuthConfig{
		ClientID:     req.ClientID,
		ClientSecret: req.ClientSecret,
		RedirectURI:  req.RedirectURI,
	}); err != nil {
		s.addLog("ERROR", fmt.Sprintf("Failed to save OAuth config: %v", err))
		http.Error(w, "Failed to save config: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update in-memory config
	newCfg := auth.OAuthConfig{
		ClientID:     req.ClientID,
		ClientSecret: req.ClientSecret,
		RedirectURI:  req.RedirectURI,
	}
	s.oauthMgr.UpdateConfig(newCfg)

	s.addLog("INFO", "OAuth config saved")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// saveOAuthConfigFile persists OAuth config to ~/Library/Application Support/kick-interval/oauth.json
func saveOAuthConfigFile(cfg auth.OAuthConfig) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configDir := homeDir + "/Library/Application Support/kick-interval"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configDir+"/oauth.json", data, 0600)
}

func (s *Server) handleChannel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Channel string `json:"channel"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Normalize and validate
	slug, err := kick.NormalizeChannel(req.Channel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := kick.ValidateSlug(slug); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Resolve channel
	s.addLog("INFO", fmt.Sprintf("Resolving channel: %s", slug))
	channel, err := s.client.GetChannel(slug)
	if err != nil {
		s.addLog("ERROR", fmt.Sprintf("Channel resolution failed: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.cfg.Channel = channel.Slug
	s.cfg.BroadcasterUserID = channel.BroadcasterUserID
	s.cfg.Save()

	s.addLog("INFO", fmt.Sprintf("Channel set: %s (ID: %d)", channel.Slug, channel.BroadcasterUserID))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"channel": channel.Slug,
		"id":      channel.BroadcasterUserID,
	})
}

func (s *Server) handleMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := kick.ValidateMessage(req.Message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.cfg.Message = req.Message
	s.cfg.Save()

	s.addLog("INFO", fmt.Sprintf("Message set: %s", req.Message))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleInterval(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Interval string `json:"interval"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	duration, err := scheduler.ParseInterval(req.Interval)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if duration < 60*time.Second {
		http.Error(w, "Interval must be at least 60 seconds", http.StatusBadRequest)
		return
	}

	s.cfg.Interval = req.Interval
	s.cfg.Save()

	s.addLog("INFO", fmt.Sprintf("Interval set: %s", req.Interval))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !s.tokenStore.HasTokens() {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	if s.cfg.Channel == "" || s.cfg.Message == "" || s.cfg.Interval == "" {
		http.Error(w, "Configuration incomplete", http.StatusBadRequest)
		return
	}

	// Parse optional body: {"send_now": true} — defaults to true
	sendNow := true
	if r.ContentLength > 0 {
		var req struct {
			SendNow *bool `json:"send_now"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if req.SendNow != nil {
			sendNow = *req.SendNow
		}
	}

	s.schedMutex.Lock()
	defer s.schedMutex.Unlock()

	if s.running {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "already_running"})
		return
	}

	interval, err := scheduler.ParseInterval(s.cfg.Interval)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var sentCount, failedCount int64

	sendFunc := func() error {
		_, err := s.client.SendChatMessage(s.cfg.BroadcasterUserID, s.cfg.Message)
		return err
	}

	onSend := func(t time.Time, err error) {
		if err != nil {
			failedCount++
			s.addLog("ERROR", fmt.Sprintf("Send failed: %v", err))
		} else {
			sentCount++
			s.addLog("INFO", "Message sent successfully")
		}
	}

	s.scheduler = scheduler.NewScheduler(interval, sendFunc, onSend, sendNow)
	s.scheduler.Start()
	s.running = true

	if sendNow {
		s.addLog("INFO", "Scheduler started — first message sent immediately")
	} else {
		s.addLog("INFO", "Scheduler started")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "started"})
}

func (s *Server) handleStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.schedMutex.Lock()
	defer s.schedMutex.Unlock()

	if !s.running {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "not_running"})
		return
	}

	if s.scheduler != nil {
		s.scheduler.Stop()
		s.scheduler = nil
	}
	s.running = false

	s.addLog("INFO", "Scheduler stopped")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "stopped"})
}

func (s *Server) handleLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Support Server-Sent Events for streaming
	if r.Header.Get("Accept") == "text/event-stream" {
		s.handleLogsSSE(w, r)
		return
	}

	// Otherwise return all logs
	s.logsMutex.Lock()
	logs := make([]LogEntry, len(s.logs))
	copy(logs, s.logs)
	s.logsMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func (s *Server) handleLogsSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Send existing logs first
	s.logsMutex.Lock()
	for _, log := range s.logs {
		data, _ := json.Marshal(log)
		fmt.Fprintf(w, "data: %s\n\n", data)
	}
	lastIndex := len(s.logs)
	s.logsMutex.Unlock()
	flusher.Flush()

	// Stream new logs
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			s.logsMutex.Lock()
			if len(s.logs) > lastIndex {
				for i := lastIndex; i < len(s.logs); i++ {
					data, _ := json.Marshal(s.logs[i])
					fmt.Fprintf(w, "data: %s\n\n", data)
				}
				lastIndex = len(s.logs)
				flusher.Flush()
			}
			s.logsMutex.Unlock()
		}
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ClearAuth   bool `json:"clear_auth"`
		ClearConfig bool `json:"clear_config"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Clear OAuth tokens (stored in keychain)
	if req.ClearAuth {
		if err := s.tokenStore.ClearTokens(); err != nil {
			s.addLog("ERROR", fmt.Sprintf("Failed to clear tokens: %v", err))
		} else {
			s.addLog("INFO", "Authentication tokens cleared")
		}
	}

	// Clear OAuth app config (client_id, client_secret — stored on disk)
	if req.ClearAuth {
		if err := clearOAuthConfigFile(); err != nil {
			s.addLog("ERROR", fmt.Sprintf("Failed to clear OAuth config file: %v", err))
		} else {
			// Reset in-memory config too
			s.oauthMgr.UpdateConfig(auth.OAuthConfig{})
			s.addLog("INFO", "OAuth app config cleared")
		}
	}

	// Clear channel/message/interval config
	if req.ClearConfig {
		s.cfg = &config.Config{}
		if err := s.cfg.Save(); err != nil {
			s.addLog("ERROR", fmt.Sprintf("Failed to clear config: %v", err))
		} else {
			s.addLog("INFO", "Configuration cleared")
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// clearOAuthConfigFile removes the oauth.json file from disk
func clearOAuthConfigFile() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := homeDir + "/Library/Application Support/kick-interval/oauth.json"
	err = os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
