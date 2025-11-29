package tui

import (
	"basicsystemmonitor/hundler"
	"fmt"
	"sort" // Import the sort package
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// MainModel holds the state of the entire TUI application.
type MainModel struct {
	cpuCh  <-chan hundler.CpuStat
	ramCh  <-chan hundler.RamStat
	diskCh <-chan hundler.DiskStat
	netCh  <-chan hundler.NetStat
	procCh <-chan []hundler.ProcessStat // New: Channel for process stats

	CpuStat       hundler.CpuStat
	RamStat       hundler.RamStat
	DiskStat      hundler.DiskStat
	NetStat       hundler.NetStat
	Processes     []hundler.ProcessStat // New: Current list of processes
	LastUpdate    time.Time

	sortBy        string // "cpu", "mem", "pid", "name"
	sortOrder     int    // 1 for ascending, -1 for descending
	ifaceName     string // New: Stores the name of the monitored network interface
	showProcesses bool   // New: Toggles process list visibility
}

// New creates a new MainModel with the given channels.
func New(cpuCh <-chan hundler.CpuStat, ramCh <-chan hundler.RamStat, diskCh <-chan hundler.DiskStat, netCh <-chan hundler.NetStat, procCh <-chan []hundler.ProcessStat, ifaceName string, showProcesses bool) MainModel {
	return MainModel{
		cpuCh:         cpuCh,
		ramCh:         ramCh,
		diskCh:        diskCh,
		netCh:         netCh,
		procCh:        procCh, // New: Store process channel
		LastUpdate:    time.Now(),
		sortBy:        "cpu", // Default sort by CPU
		sortOrder:     -1,    // Default descending
		ifaceName:     ifaceName, // New: Store network interface name
		showProcesses: showProcesses, // New: Store process list visibility
	}
}

// Msg types for updating the model
type cpuMsg hundler.CpuStat
type ramMsg hundler.RamStat
type diskMsg hundler.DiskStat
type netMsg hundler.NetStat
type processMsg []hundler.ProcessStat // New: Message type for process stats
type tickMsg time.Time

// waitForActivity is a command that waits for activity on any of the stat channels.
func (m *MainModel) waitForActivity() tea.Cmd {
	return func() tea.Msg {
		select {
		case cpu := <-m.cpuCh:
			return cpuMsg(cpu)
		case ram := <-m.ramCh:
			return ramMsg(ram)
		case disk := <-m.diskCh:
			return diskMsg(disk)
		case net := <-m.netCh:
			return netMsg(net)
		case procs := <-m.procCh: // New: Listen for process updates
			return processMsg(procs)
		}
	}
}

func tickCommand(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Init initializes the model.
func (m MainModel) Init() tea.Cmd {
	return tea.Batch(
		m.waitForActivity(),
		tickCommand(time.Second),
	)
}

// sortProcesses sorts the process list based on the current sortBy and sortOrder.
func (m *MainModel) sortProcesses() {
	sort.Slice(m.Processes, func(i, j int) bool {
		var less bool
		switch m.sortBy {
		case "pid":
			less = m.Processes[i].Pid < m.Processes[j].Pid
		case "name":
			less = m.Processes[i].Name < m.Processes[j].Name
		case "cpu":
			less = m.Processes[i].CPUPercent < m.Processes[j].CPUPercent
		case "mem":
			less = m.Processes[i].MemoryBytes < m.Processes[j].MemoryBytes
		default:
			less = m.Processes[i].CPUPercent < m.Processes[j].CPUPercent
		}

		if m.sortOrder == -1 {
			return !less
		}
		return less
	})
}

// Update handles messages and updates the model accordingly.
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "c": // Sort by CPU
			if m.sortBy == "cpu" {
				m.sortOrder *= -1 // Toggle order
			} else {
				m.sortBy = "cpu"
				m.sortOrder = -1 // Default descending for CPU
			}
			m.sortProcesses() // Sort immediately after key press
		case "m": // Sort by Memory
			if m.sortBy == "mem" {
				m.sortOrder *= -1 // Toggle order
			} else {
				m.sortBy = "mem"
				m.sortOrder = -1 // Default descending for Memory
			}
			m.sortProcesses() // Sort immediately after key press
		case "p": // Sort by PID
			if m.sortBy == "pid" {
				m.sortOrder *= -1 // Toggle order
			} else {
				m.sortBy = "pid"
				m.sortOrder = 1 // Default ascending for PID
			}
			m.sortProcesses() // Sort immediately after key press
		case "n": // Sort by Name
			if m.sortBy == "name" {
				m.sortOrder *= -1 // Toggle order
			} else {
				m.sortBy = "name"
				m.sortOrder = 1 // Default ascending for Name
			}
			m.sortProcesses() // Sort immediately after key press
		}
	case cpuMsg:
		m.CpuStat = hundler.CpuStat(msg)
		return m, m.waitForActivity()
	case ramMsg:
		m.RamStat = hundler.RamStat(msg)
		return m, m.waitForActivity()
	case diskMsg:
		m.DiskStat = hundler.DiskStat(msg)
		return m, m.waitForActivity()
	case netMsg:
		m.NetStat = hundler.NetStat(msg)
		return m, m.waitForActivity()
	case processMsg: // New: Handle process updates
		m.Processes = processMsg(msg)
		m.sortProcesses() // Sort after receiving new data
		return m, m.waitForActivity()
	case tickMsg:
		m.LastUpdate = time.Time(msg)
		return m, tickCommand(time.Second)
	}
	return m, nil
}

// View renders the UI.
func (m MainModel) View() string {
	s := fmt.Sprintf("Basic System Monitor — %s\n\n", m.LastUpdate.Format(time.RFC1123))

	s += fmt.Sprintf("CPU:           %6.2f%% \n\n", m.CpuStat.Percent)
	s += fmt.Sprintf("RAM:           %8s / %8s (%6.2f%%)\n\n", ByteCountSI(m.RamStat.Used), ByteCountSI(m.RamStat.Total), m.RamStat.UsedPercent)
	s += fmt.Sprintf("Disk (/):      %8s (%6.2f%%) \n\n", ByteCountSI(m.DiskStat.Used), m.DiskStat.UsedPercent)
	netInfo := "Network:"
	if m.ifaceName != "" {
		netInfo += fmt.Sprintf(" (%s)", m.ifaceName)
	}
	s += fmt.Sprintf("%-15s ↑ %8s/s   ↓ %8s/s\n\n", netInfo, ByteCountSI(uint64(m.NetStat.BytesSentPerSec)), ByteCountSI(uint64(m.NetStat.BytesRecvPerSec)))

	if m.showProcesses {
		s += "Processes:\n"
		s += fmt.Sprintf("%-8s %-30s %-8s %-8s\n", "PID", "NAME", "CPU%", "MEM")
		for i, p := range m.Processes {
			if i >= 10 { // Limit to 10 processes for brevity for now
				break
			}
			s += fmt.Sprintf("%-8d %-30s %-7.2f%% %s\n", p.Pid, p.Name, p.CPUPercent, ByteCountSI(p.MemoryBytes))
		}
	}

	s += "\nPress 'q' or 'ctrl+c' to quit.\n"
	return s
}
