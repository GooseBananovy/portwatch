# portwatch

[![Go Version](https://img.shields.io/github/go-mod/go-version/goosebananovy/portwatch)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**portwatch** is a lightweight Linux server monitoring agent with a CLI client. Written in Go, deployed via Docker. The agent collects real‑time system metrics and exposes them over a simple HTTP API, while the client lets you query the data from any machine.

## Architecture

```
cmd/
  agent/        ← HTTP server, runs on the remote server
  cli/          ← CLI client, runs on your machine
internal/
  handler/      ← HTTP request handlers
  model/        ← metric types (CPU, RAM, Disk, Network, etc.)
  provider/
    linux/      ← reads metrics from /proc and Docker API
```

## Collected Metrics

| Category | Metrics |
|----------|---------|
| CPU      | Core count, per‑core load (%), total load (%) |
| RAM      | Total, used, available (bytes) |
| Disk     | Per partition: path, total, used, free (bytes) |
| Uptime   | System uptime in seconds |
| Network  | Incoming/outgoing bytes per second, active TCP connections |
| Processes| Top 30 processes by CPU: name, PID, CPU%, RAM% |
| Docker   | Per container: name, status, CPU%, RAM used / total RAM |

## API Endpoints

All endpoints return JSON.

| Method | Path | Description |
|--------|------|-------------|
| GET | `/portwatch/health` | Health check |
| GET | `/portwatch/uptime` | System uptime |
| GET | `/portwatch/cpu` | CPU metrics |
| GET | `/portwatch/ram` | RAM metrics |
| GET | `/portwatch/disk` | Disk metrics |
| GET | `/portwatch/sys` | All system metrics (CPU, RAM, Disk, Uptime) |
| GET | `/portwatch/network` | Network metrics |
| GET | `/portwatch/procs` | Top 30 processes |
| GET | `/portwatch/docker` | Docker container info |
| GET | `/portwatch/all` | Everything at once |

## Quick Start (Local Run)

### Prerequisites
- Go 1.21+ (if building from source)
- Linux system (agent requires `/proc`)

### Run the Agent
```bash
git clone https://github.com/goosebananovy/portwatch
cd portwatch

# Start the agent (listens on port 9090 by default)
go run cmd/agent/main.go
```

### Run the CLI
In another terminal:
```bash
# Point to the running agent
export PORTWATCH_HOST=http://localhost:9090

# Fetch all metrics
go run cmd/cli/main.go all
```

Available CLI commands: `health`, `uptime`, `cpu`, `ram`, `disk`, `sys`, `network`, `procs`, `docker`, `all`.

## Install the CLI as a Binary

You can install the CLI client globally using `go install`:
```bash
go install github.com/goosebananovy/portwatch/cmd/cli@latest
```
After that, the `portwatch` command will be available system‑wide (provided `$GOPATH/bin` is in your `PATH`).

## Deploy the Agent on a Server with Docker

### Requirements
- Docker and Docker Compose on the target Linux machine
- `/proc` filesystem (available by default)

### Steps
```bash
git clone https://github.com/goosebananovy/portwatch
cd portwatch

# Run the agent in the background
docker compose up --build -d
```

After startup, the agent will be reachable on port `90` of your server (mapped to internal port 9090).

### Environment Variables for Docker

| Variable | Default | Description |
|----------|---------|-------------|
| `PROC_BASE` | `/proc` | Path to the proc filesystem |
| `PID1_PROC_BASE` | `$PROC_BASE` | Path to the proc namespace of PID 1 (used to correctly read network and disk stats when running inside a container) |

Example `docker-compose.yml` (already included in the repo):
```yaml
services:
  app:
    build:
      context: .
      dockerfile: Dockerfile.agent
    ports:
      - "90:9090"
    environment:
      - PROC_BASE=/host/proc
      - PID1_PROC_BASE=/host/proc/1
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - /proc:/host/proc:ro
```

## Using the CLI

Once the CLI is installed (or via `go run`), you can query the agent:

```bash
# Set the remote server address (or keep localhost for local runs)
export PORTWATCH_HOST=http://your-server:90

# Fetch all metrics
portwatch all

# CPU metrics only
portwatch cpu

# Docker container info
portwatch docker

# Top processes
portwatch procs
```

## In proress

- Beautiful TUI with Bubble Tea
- In‑terminal metric charts
- Sensor metrics (CPU temperature, fan speed)
- Caddy reverse proxy for port‑free access
- Unit tests
- Historical metrics with PostgreSQL

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
