package hundler

import (
	"context"
	"time"

	"github.com/shirou/gopsutil/v4/disk"
)

type DiskStat struct {
	Path        string
	Total       uint64
	Used        uint64
	UsedPercent float64
}

// StartDiskMonitor streams usage for a given path (e.g. "/") periodically.
func StartDiskMonitor(ctx context.Context, interval time.Duration, path string) <-chan DiskStat {
	ch := make(chan DiskStat)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			u, err := disk.Usage(path)
			var s DiskStat
			if err == nil {
				s = DiskStat{Path: path, Total: u.Total, Used: u.Used, UsedPercent: u.UsedPercent}
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
