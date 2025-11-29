package hundler

import (
	"context"
	"log"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

// ProcessStat holds periodic process information
type ProcessStat struct {
	Pid         int32
	Name        string
	CPUPercent  float64
	MemoryBytes uint64
}

// StartProcessMonitor starts a goroutine that periodically sends a list of process stats to the returned channel.
// The channel is closed when the provided context is cancelled.
func StartProcessMonitor(ctx context.Context, interval time.Duration) <-chan []ProcessStat {
	ch := make(chan []ProcessStat)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			procs, err := process.Processes()
			if err != nil {
				log.Printf("Error getting processes: %v", err)
				select {
				case <-ticker.C:
					continue
				case <-ctx.Done():
					return
				}
			}

			var stats []ProcessStat
			for _, p := range procs {
				name, err := p.NameWithContext(ctx)
				if err != nil {
					continue
				}
				cpuPercent, err := p.CPUPercentWithContext(ctx)
				if err != nil {
					continue
				}
				memInfo, err := p.MemoryInfoWithContext(ctx)
				if err != nil {
					continue
				}

				stats = append(stats, ProcessStat{
					Pid:         p.Pid,
					Name:        name,
					CPUPercent:  cpuPercent,
					MemoryBytes: memInfo.RSS, // Resident Set Size
				})
			}

			select {
			case ch <- stats:
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