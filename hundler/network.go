package hundler

import (
	"context"
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
func StartNetworkMonitor(ctx context.Context, interval time.Duration) <-chan NetStat {
	ch := make(chan NetStat)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		var prevSent uint64
		var prevRecv uint64
		// initial read
		if counters, err := net.IOCounters(false); err == nil && len(counters) > 0 {
			prevSent = counters[0].BytesSent
			prevRecv = counters[0].BytesRecv
		}

		for {
			counters, err := net.IOCounters(false)
			var s NetStat
			if err == nil && len(counters) > 0 {
				cur := counters[0]
				sent := cur.BytesSent
				recv := cur.BytesRecv
				elapsed := interval.Seconds()
				if prevSent > 0 {
					s.BytesSentPerSec = float64(sent-prevSent) / elapsed
				}
				if prevRecv > 0 {
					s.BytesRecvPerSec = float64(recv-prevRecv) / elapsed
				}
				s.TotalBytesSent = sent
				s.TotalBytesRecv = recv
				prevSent = sent
				prevRecv = recv
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
