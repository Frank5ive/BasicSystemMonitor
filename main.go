package main

import (
	"basicsystemmonitor/hundler"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// start monitors
	cpuCh := hundler.StartCpuMonitor(ctx, 1*time.Second)
	ramCh := hundler.StartRamMonitor(ctx, 1*time.Second)
	diskCh := hundler.StartDiskMonitor(ctx, 2*time.Second, "/")
	netCh := hundler.StartNetworkMonitor(ctx, 1*time.Second)

	var cpuS hundler.CpuStat
	var ramS hundler.RamStat
	var diskS hundler.DiskStat
	var netS hundler.NetStat

	refresh := time.NewTicker(500 * time.Millisecond)
	defer refresh.Stop()

	// main loop: update latest stats and periodically redraw
	for {
		select {
		case s, ok := <-cpuCh:
			if ok {
				cpuS = s
			}
		case s, ok := <-ramCh:
			if ok {
				ramS = s
			}
		case s, ok := <-diskCh:
			if ok {
				diskS = s
			}
		case s, ok := <-netCh:
			if ok {
				netS = s
			}
		case <-refresh.C:
			printScreen(cpuS, ramS, diskS, netS)
		case <-sigs:
			cancel()
			// let monitors clean up
			time.Sleep(200 * time.Millisecond)
			fmt.Println("\nExiting...")
			return
		case <-ctx.Done():
			return
		}
	}
}

func printScreen(cpuS hundler.CpuStat, ramS hundler.RamStat, diskS hundler.DiskStat, netS hundler.NetStat) {
	// clear screen + move cursor home
	fmt.Print("\033[2J\033[H")
	fmt.Printf("Basic System Monitor — %s\n\n", time.Now().Format(time.RFC1123))

	fmt.Printf("CPU: %.2f%%\n", cpuS.Percent)
	fmt.Println()
	fmt.Printf("RAM: Used %s / %s (%.2f%%)\n", byteCountSI(ramS.Used), byteCountSI(ramS.Total), ramS.UsedPercent)
	fmt.Println()
	fmt.Printf("Disk (%s): Used %s (%.2f%%)\n", diskS.Path, byteCountSI(diskS.Used), diskS.UsedPercent)
	fmt.Println()
	fmt.Printf("Network: ↑ %s/s  ↓ %s/s\n", byteCountSI(uint64(netS.BytesSentPerSec)), byteCountSI(uint64(netS.BytesRecvPerSec)))
}

// byteCountSI formats bytes in SI units (kB=1000)
func byteCountSI(b uint64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
