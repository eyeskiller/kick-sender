package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Scheduler manages periodic message sending
type Scheduler struct {
	interval    time.Duration
	sendFunc    func() error
	onSend      func(time.Time, error)
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	sendNow     bool
	mu          sync.Mutex
	sending     bool
	nextSendMu  sync.RWMutex
	nextSendAt  time.Time
}

// NewScheduler creates a new scheduler
func NewScheduler(interval time.Duration, sendFunc func() error, onSend func(time.Time, error), sendNow bool) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		interval: interval,
		sendFunc: sendFunc,
		onSend:   onSend,
		ctx:      ctx,
		cancel:   cancel,
		sendNow:  sendNow,
	}
}

// Start begins the scheduler
func (s *Scheduler) Start() {
	s.wg.Add(1)
	go s.run()
}

func (s *Scheduler) run() {
	defer s.wg.Done()

	// Compute first target time
	s.nextSendMu.Lock()
	s.nextSendAt = time.Now().Add(s.interval)
	s.nextSendMu.Unlock()

	// Send immediately if requested
	if s.sendNow {
		s.send()
	}

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case t := <-ticker.C:
			// The ticker just fired — the next send is one interval ahead of this tick.
			s.nextSendMu.Lock()
			s.nextSendAt = t.Add(s.interval)
			s.nextSendMu.Unlock()

			s.send()
		}
	}
}

func (s *Scheduler) send() {
	s.mu.Lock()
	if s.sending {
		s.mu.Unlock()
		return
	}
	s.sending = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.sending = false
		s.mu.Unlock()
	}()

	err := s.sendFunc()
	if s.onSend != nil {
		s.onSend(time.Now(), err)
	}
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.cancel()
	s.wg.Wait()
}

// NextSend returns the time of the next scheduled send.
func (s *Scheduler) NextSend() time.Time {
	s.nextSendMu.RLock()
	defer s.nextSendMu.RUnlock()
	if s.nextSendAt.IsZero() {
		return time.Time{}
	}
	return s.nextSendAt
}

// FormatInterval formats a duration as a human-readable string
func FormatInterval(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	return fmt.Sprintf("%dh", int(d.Hours()))
}

// ParseInterval parses an interval string like "10m", "1h", "60s"
func ParseInterval(s string) (time.Duration, error) {
	var value int
	var unit string

	_, err := fmt.Sscanf(s, "%d%s", &value, &unit)
	if err != nil {
		return 0, fmt.Errorf("invalid interval format: %s", s)
	}

	switch unit {
	case "s":
		return time.Duration(value) * time.Second, nil
	case "m":
		return time.Duration(value) * time.Minute, nil
	case "h":
		return time.Duration(value) * time.Hour, nil
	default:
		return 0, fmt.Errorf("invalid interval unit: %s (use s, m, or h)", unit)
	}
}
