package hundler

import (
	"context"
	"testing"
	"time"
)

func TestStartDiskMonitor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interval := 100 * time.Millisecond
	path := "/" // Test root path
	diskCh := StartDiskMonitor(ctx, interval, path)

	// Test if at least one value is received
	select {
	case stat := <-diskCh:
		if stat.Total == 0 {
			t.Errorf("DiskStat.Total should not be 0")
		}
		if stat.UsedPercent < 0 || stat.UsedPercent > 100 {
			t.Errorf("DiskStat.UsedPercent out of range: %f", stat.UsedPercent)
		}
		if stat.Path != path {
			t.Errorf("DiskStat.Path mismatch: got %s, want %s", stat.Path, path)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Timeout waiting for Disk stats")
	}

	// Test if monitoring stops after context cancellation
	cancel()
	select {
	case _, ok := <-diskCh:
		if ok {
			t.Error("Disk channel should be closed after context cancellation")
		}
	case <-time.After(500 * time.Millisecond):
		// This is acceptable
	}
}