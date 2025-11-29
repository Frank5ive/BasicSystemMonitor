package hundler

import (
	"context"
	"testing"
	"time"
)

func TestStartRamMonitor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interval := 100 * time.Millisecond
	ramCh := StartRamMonitor(ctx, interval)

	// Test if at least one value is received
	select {
	case stat := <-ramCh:
		if stat.Total == 0 {
			t.Errorf("RamStat.Total should not be 0")
		}
		if stat.UsedPercent < 0 || stat.UsedPercent > 100 {
			t.Errorf("RamStat.UsedPercent out of range: %f", stat.UsedPercent)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Timeout waiting for RAM stats")
	}

	// Test if monitoring stops after context cancellation
	cancel()
	select {
	case _, ok := <-ramCh:
		if ok {
			t.Error("RAM channel should be closed after context cancellation")
		}
	case <-time.After(500 * time.Millisecond):
		// This is acceptable
	}
}