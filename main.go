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

	// initialize screen once (clear, hide cursor, draw labels)
	initScreen()
	defer restoreTerminal()

	// main loop: update latest stats and periodically redraw only fields
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
			updateScreen(cpuS, ramS, diskS, netS)
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

// Terminal helpers: draw static layout once, then update fields in-place.
func initScreen() {
	// clear screen and hide cursor
	fmt.Print("\033[2J")
	fmt.Print("\033[?25l")
	// draw static labels
	move(1, 1)
	fmt.Printf("Basic System Monitor — %s", time.Now().Format(time.RFC1123))

	move(3, 1)
	fmt.Print("CPU:           ")

	move(5, 1)
	fmt.Print("RAM:           Used / Total (%%)")

	move(7, 1)
	fmt.Print("Disk (/):      Used (%%)")

	move(9, 1)
	fmt.Print("Network:       ↑ /s   ↓ /s")
}

func restoreTerminal() {
	// show cursor and move to next line
	fmt.Print("\033[?25h")
	fmt.Print("\n")
}

func updateScreen(cpuS hundler.CpuStat, ramS hundler.RamStat, diskS hundler.DiskStat, netS hundler.NetStat) {
	// update timestamp
	move(1, 1)
	fmt.Printf("Basic System Monitor — %s", time.Now().Format(time.RFC1123))

	// print CPU at col 16 (after label), fixed width
	move(3, 16)
	fmt.Printf("%6.2f%%     ", cpuS.Percent)

	// RAM: Used / Total (percent)
	move(5, 16)
	fmt.Printf("%8s / %8s (%6.2f%%)", byteCountSI(ramS.Used), byteCountSI(ramS.Total), ramS.UsedPercent)

	// Disk
	move(7, 16)
	fmt.Printf("%8s (%6.2f%%)   ", byteCountSI(diskS.Used), diskS.UsedPercent)

	// Network
	move(9, 16)
	fmt.Printf("%8s   %8s   ", byteCountSI(uint64(netS.BytesSentPerSec)), byteCountSI(uint64(netS.BytesRecvPerSec)))
}

// move cursor to (row, col)
func move(row, col int) {
	fmt.Printf("\033[%d;%dH", row, col)
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
