# Basic System Monitor

A terminal-based system monitor written in Go, displaying CPU, RAM, disk, and network usage, along with a sortable process list.

## Features

- **Real-time Monitoring:** Concurrent monitors for CPU, RAM, Disk, and Network usage.
- **Interactive Process List:** View a list of running processes with their PID, Name, CPU usage, and Memory usage.
- **Sortable Processes:** Sort the process list by PID, Name, CPU, or Memory by pressing 'p', 'n', 'c', or 'm' respectively.
- **Configurable:** Customize refresh intervals and other settings via `config.yaml` or command-line flags.
- **Network Interface Selection:** Monitor a specific network interface.
- **Process List Visibility:** Show or hide the process list with a command-line flag.
- **Docker Support:** A multi-stage `Dockerfile` is provided for building a small, efficient container image.

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/your-username/BasicSystemMonitor.git
cd BasicSystemMonitor

# Build the binary
go build -o basic-system-monitor ./

# Run the monitor
./basic-system-monitor
```

### Docker

```bash
# Build the Docker image
docker build -t basic-system-monitor:latest .

# Run the monitor
docker run --rm -it --name monitor basic-system-monitor:latest
```

## Usage

```bash
./basic-system-monitor [flags]
```

### Command-Line Flags

| Flag        | Description                                       | Default     |
|-------------|---------------------------------------------------|-------------|
| `-c`        | Path to configuration file                        | `config.yaml` |
| `-i`        | Refresh interval (e.g., 1s, 500ms)                | `1s`        |
| `-d`        | Disk path to monitor (e.g., /var, C:\\)             | `/`         |
| `-iface`    | Network interface to monitor (e.g., eth0, en0)    | (all)       |
| `-p`        | Show process list                                 | `false`     |

### Configuration

The application can be configured via a `config.yaml` file. Command-line flags will override the values in the config file.

**Example `config.yaml`:**
```yaml
refreshInterval: 2s
diskPath: /home
```

## Interactive Controls

- `q`, `ctrl+c`: Quit the application.
- `c`: Sort processes by CPU usage.
- `m`: Sort processes by Memory usage.
- `p`: Sort processes by PID.
- `n`: Sort processes by Name.

## License

MIT
