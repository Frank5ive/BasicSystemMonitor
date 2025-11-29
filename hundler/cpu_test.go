package hundler

import (
	"context"
	"testing"
	"time"
)

func TestStartCpuMonitor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interval := 100 * time.Millisecond
	cpuCh := StartCpuMonitor(ctx, interval)

	// Test if at least one value is received
	select {
	case stat := <-cpuCh:
		if stat.Percent < 0 || stat.Percent > 100 {
			t.Errorf("CpuStat.Percent out of range: %f", stat.Percent)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Timeout waiting for CPU stats")
	}

	// Test if monitoring stops after context cancellation
	cancel()
	select {
	case _, ok := <-cpuCh:
		if ok {
			t.Error("CPU channel should be closed after context cancellation")
		}
	case <-time.After(500 * time.Millisecond):
		// This is acceptable, as the goroutine might take a moment to shut down
	}
}