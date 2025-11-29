package hundler

import (
	"context"
	"log"
	"time"

	"github.com/shirou/gopsutil/v4/net"
)

type NetStat struct {
	BytesSentPerSec float64
	BytesRecvPerSec float64
	TotalBytesSent  uint64
	TotalBytesRecv  uint64
}

// StartNetworkMonitor streams network byte rates (per-second) computed from IO counters.
func StartNetworkMonitor(ctx context.Context, interval time.Duration, ifaceName string) <-chan NetStat {
	ch := make(chan NetStat)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		var prevSent uint64
		var prevRecv uint64

		// Initial read to set prevSent and prevRecv
		if ifaceName != "" {
			if counters, err := net.IOCounters(true); err == nil {
				for _, c := range counters {
					if c.Name == ifaceName {
						prevSent = c.BytesSent
						prevRecv = c.BytesRecv
						break
					}
				}
			}
		} else {
			if counters, err := net.IOCounters(false); err == nil && len(counters) > 0 {
				prevSent = counters[0].BytesSent
				prevRecv = counters[0].BytesRecv
			}
		}

		for {
			var curSent uint64
			var curRecv uint64
			var err error

			if ifaceName != "" {
				var found bool
				var counters []net.IOCountersStat
				counters, err = net.IOCounters(true)
				if err == nil {
					for _, c := range counters {
						if c.Name == ifaceName {
							curSent = c.BytesSent
							curRecv = c.BytesRecv
							found = true
							break
						}
					}
				}
				if !found && err == nil {
					// Interface not found, perhaps it was removed or never existed
					// Log this or send an error state if needed
					log.Printf("Network interface '%s' not found.", ifaceName)
				}
			} else {
				var counters []net.IOCountersStat
				counters, err = net.IOCounters(false)
				if err == nil && len(counters) > 0 {
					curSent = counters[0].BytesSent
					curRecv = counters[0].BytesRecv
				}
			}

			var s NetStat
			if err == nil {
				elapsed := interval.Seconds()
				if prevSent > 0 {
					s.BytesSentPerSec = float64(curSent-prevSent) / elapsed
				}
				if prevRecv > 0 {
					s.BytesRecvPerSec = float64(curRecv-prevRecv) / elapsed
				}
				s.TotalBytesSent = curSent
				s.TotalBytesRecv = curRecv
				prevSent = curSent
				prevRecv = curRecv
			}

			select {
			case ch <- s:
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
