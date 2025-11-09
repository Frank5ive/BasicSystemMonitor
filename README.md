# BasicSystemMonitor

A small Go-based terminal system monitor that streams CPU, RAM, Disk and Network usage in real time using `gopsutil`.

## Features

- Concurrent monitors for CPU, RAM, Disk and Network (implemented in `hundler/`).
- Simple ANSI terminal UI refreshed periodically (no external UI libs required).
- Ready-to-build Dockerfile for containerized execution.

## Build & Run (local)

From the project root:

```bash
# fetch dependencies and build
go mod tidy
go build -o monitor ./

# run
./monitor
```

You can also run with `go run .` for development.

## Docker

Build the Docker image and run it:

```bash
# build
docker build -t basicsystemmonitor:latest .

# run (interactive tty so you can see the terminal UI)
docker run --rm -it --name monitor basicsystemmonitor:latest
```

Notes:

- The monitor reads system stats via `gopsutil`. When running in a container, container resource limits and the runtime environment affect reported values.

## Files of interest

- `main.go` — program entry; starts monitors and redraws terminal UI.
- `hundler/` — monitoring implementations:
  - `cpu.go` — CPU monitor (`StartCpuMonitor`).
  - `ram.go` — RAM monitor (`StartRamMonitor`).
  - `disk.go` — Disk monitor (`StartDiskMonitor`).
  - `network.go` — Network monitor (`StartNetworkMonitor`).

## Contributing

- Open an issue or PR to add features (per-core graphs, colored UI, CLI flags, etc.).

## License

MIT
