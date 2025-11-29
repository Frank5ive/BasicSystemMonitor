# Basic System Monitor

```
   ____    _    ____   ___   ____  _  _   _   _  ____  _   _ _____
  | __ )  / \  / ___| / _ \ / ___|| || | | | | |/ ___|| | | | ____|
  |  _ \ / _ \ \___ \| | | | |    | || |_| | | | |    | |_| |  _|
  | |_) / ___ \ ___) | |_| | |___ |__   _| |_| | |___ |  _  | |___
  |____/_/   \_\____/ \___/ \____|   |_|  \___/ \____||_| |_|_____|
```

**A modern, lightweight, and highly configurable system monitoring tool for your terminal.**

[![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

---

Basic System Monitor is a modern, lightweight, and highly configurable system monitoring tool for your terminal. Written in Go and built with the `bubbletea` TUI framework, it provides real-time insights into your system's CPU, RAM, disk, and network usage, along with a sortable process list. With a focus on performance and ease of use, Basic System Monitor is a great tool for developers and system administrators who want a quick and easy way to keep an eye on their system's health. The application is highly configurable through a `config.yaml` file and command-line flags, allowing you to tailor it to your specific needs.

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
| `-d`        | Disk path to monitor (e.g., /var, C:\)             | `/`         |
| `-iface`    | Network interface to monitor (e.g., eth0, en0)    | (all)       |
| `-p`        | Show process list                                 | `false`     |
| `-proc-interval`| Process list refresh interval (e.g., 3s, 5s) | `3s`        |


### Configuration

The application can be configured via a `config.yaml` file. Command-line flags will override the values in the config file.

**Example `config.yaml`:**
```yaml
refreshInterval: 2s
diskPath: /home
processRefreshInterval: 5s
```

## Interactive Controls

- `q`, `ctrl+c`: Quit the application.
- `c`: Sort processes by CPU usage.
- `m`: Sort processes by Memory usage.
- `p`: Sort processes by PID.
- `n`: Sort processes by Name.

## Contributing

Contributions are welcome! Please feel free to open an issue or submit a pull request.

1.  Fork the repository.
2.  Create your feature branch (`git checkout -b feature/AmazingFeature`).
3.  Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4.  Push to the branch (`git push origin feature/AmazingFeature`).
5.  Open a pull request.

Please make sure to update tests as appropriate.

## Code of Conduct

Please note that this project is released with a Contributor Code of Conduct. By participating in this project you agree to abide by its terms. Please see [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) for more information.

## License

Distributed under the MIT License. See `LICENSE` for more information.