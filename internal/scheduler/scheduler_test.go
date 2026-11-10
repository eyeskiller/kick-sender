package scheduler

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestParseInterval(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Duration
		wantErr bool
	}{
		{"seconds", "60s", 60 * time.Second, false},
		{"minutes", "10m", 10 * time.Minute, false},
		{"hours", "1h", 1 * time.Hour, false},
		{"invalid format", "invalid", 0, true},
		{"invalid unit", "10x", 0, true},
		{"empty", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseInterval(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInterval(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseInterval(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestFormatInterval(t *testing.T) {
	tests := []struct {
		name  string
		input time.Duration
		want  string
	}{
		{"seconds", 30 * time.Second, "30s"},
		{"minutes", 10 * time.Minute, "10m"},
		{"hours", 2 * time.Hour, "2h"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatInterval(tt.input)
			if got != tt.want {
				t.Errorf("FormatInterval(%v) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSchedulerCancellation(t *testing.T) {
	var callCount int64

	sendFunc := func() error {
		atomic.AddInt64(&callCount, 1)
		return nil
	}

	onSend := func(t time.Time, err error) {}

	// Create scheduler with short interval
	sched := NewScheduler(100*time.Millisecond, sendFunc, onSend, true)
	sched.Start()

	// Let it run for a bit
	time.Sleep(350 * time.Millisecond)

	sched.Stop()

	count := atomic.LoadInt64(&callCount)
	// Should have sent at least 2 messages (immediate + 2 ticks)
	if count < 2 {
		t.Errorf("expected at least 2 sends, got %d", count)
	}
}

func TestSchedulerNoOverlap(t *testing.T) {
	var sending int64
	var maxConcurrent int64

	sendFunc := func() error {
		current := atomic.AddInt64(&sending, 1)
		for {
			old := atomic.LoadInt64(&maxConcurrent)
			if current <= old || atomic.CompareAndSwapInt64(&maxConcurrent, old, current) {
				break
			}
		}
		time.Sleep(200 * time.Millisecond) // Simulate slow send
		atomic.AddInt64(&sending, -1)
		return nil
	}

	onSend := func(t time.Time, err error) {}

	// Create scheduler with very short interval
	sched := NewScheduler(50*time.Millisecond, sendFunc, onSend, true)
	sched.Start()

	time.Sleep(500 * time.Millisecond)
	sched.Stop()

	max := atomic.LoadInt64(&maxConcurrent)
	if max > 1 {
		t.Errorf("expected no overlapping sends, got max concurrent = %d", max)
	}
}

func TestSchedulerGracefulShutdown(t *testing.T) {
	sendFunc := func() error {
		return nil
	}

	onSend := func(t time.Time, err error) {}

	sched := NewScheduler(1*time.Second, sendFunc, onSend, false)
	sched.Start()

	// Should stop cleanly
	done := make(chan struct{})
	go func() {
		sched.Stop()
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(2 * time.Second):
		t.Error("scheduler did not stop within timeout")
	}
}
