package hundler

import (
	"context"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

// CpuStat holds periodic CPU usage information
type CpuStat struct {
	Percent float64
}

// StartCpuMonitor starts a goroutine that periodically sends CPU usage percent to the returned channel.
// The channel is closed when the provided context is cancelled.
func StartCpuMonitor(ctx context.Context, interval time.Duration) <-chan CpuStat {
	ch := make(chan CpuStat)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			// sample
			percents, err := cpu.Percent(0, false)
			var p float64
			if err == nil && len(percents) > 0 {
				p = percents[0]
			}

			select {
			case ch <- CpuStat{Percent: p}:
			case <-ctx.Done():
				return
			}

			select {
			case <-ticker.C:
				continue
			case <-ctx.Done():
				return
			}
		}
	}()
	return ch
}
