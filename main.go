package main

import (
	"basicsystemmonitor/hundler"
	"basicsystemmonitor/tui"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var configPath string
	var refreshIntervalStr string
	var diskPathStr string
	var ifaceName string
	var showProcesses bool // New: for process list visibility

	flag.StringVar(&configPath, "c", "config.yaml", "Path to configuration file")
	flag.StringVar(&refreshIntervalStr, "i", "", "Refresh interval (e.g., 1s, 500ms)")
	flag.StringVar(&diskPathStr, "d", "", "Disk path to monitor (e.g., /var, C:\\)")
	flag.StringVar(&ifaceName, "iface", "", "Network interface to monitor (e.g., eth0, en0)")
	flag.BoolVar(&showProcesses, "p", false, "Show process list") // New flag
	flag.Parse()

	config, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Override config values with command-line flags if provided
	if refreshIntervalStr != "" {
		config.RefreshInterval = refreshIntervalStr
	}
	if diskPathStr != "" {
		config.DiskPath = diskPathStr
	}

	refreshInterval, err := config.GetRefreshInterval()
	if err != nil {
		log.Fatalf("Error parsing refresh interval: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start monitors and get their channels
	cpuCh := hundler.StartCpuMonitor(ctx, refreshInterval)
	ramCh := hundler.StartRamMonitor(ctx, refreshInterval)
	diskCh := hundler.StartDiskMonitor(ctx, 2*refreshInterval, config.DiskPath)
	netCh := hundler.StartNetworkMonitor(ctx, refreshInterval, ifaceName) // Pass ifaceName
	procCh := hundler.StartProcessMonitor(ctx, refreshInterval) // Start process monitor

	// Initialize the Bubble Tea model with the channels
	initialModel := tui.New(cpuCh, ramCh, diskCh, netCh, procCh, ifaceName, showProcesses)

	// Start the Bubble Tea program
	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
