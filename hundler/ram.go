package hundler

import (
	"context"
	"time"

	"github.com/shirou/gopsutil/v4/mem"
)

type RamStat struct {
	Total       uint64
	Used        uint64
	UsedPercent float64
}

func StartRamMonitor(ctx context.Context, interval time.Duration) <-chan RamStat {
	ch := make(chan RamStat)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			v, err := mem.VirtualMemory()
			var s RamStat
			if err == nil {
				s = RamStat{Total: v.Total, Used: v.Used, UsedPercent: v.UsedPercent}
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
