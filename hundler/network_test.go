package hundler

import (
	"context"
	"testing"
	"time"
)

func TestStartNetworkMonitor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interval := 100 * time.Millisecond
	netCh := StartNetworkMonitor(ctx, interval)

	// Test if at least one value is received
	select {
	case stat := <-netCh:
		// We can't assert specific values for BytesSentPerSec/BytesRecvPerSec
		// as they depend on network activity, but we can check if they are non-negative.
		if stat.BytesSentPerSec < 0 {
			t.Errorf("NetStat.BytesSentPerSec should not be negative: %f", stat.BytesSentPerSec)
		}
		if stat.BytesRecvPerSec < 0 {
			t.Errorf("NetStat.BytesRecvPerSec should not be negative: %f", stat.BytesRecvPerSec)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Timeout waiting for Network stats")
	}

	// Test if monitoring stops after context cancellation
	cancel()
	select {
	case _, ok := <-netCh:
		if ok {
			t.Error("Network channel should be closed after context cancellation")
		}
	case <-time.After(500 * time.Millisecond):
		// This is acceptable
	}
}