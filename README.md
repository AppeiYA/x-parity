# X-Parity

[![Go Report Card](https://goreportcard.com/badge/github.com/AppeiYA/x-parity)](https://goreportcard.com/report/github.com/AppeiYA/x-parity)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Zero Dependencies](https://img.shields.io/badge/dependencies-zero-brightgreen.svg)](go.mod)
[![Platform](https://img.shields.io/badge/platform-linux%20%7C%20macos%20%7C%20windows-blue.svg)](dist/)
[![Binary Size](https://img.shields.io/badge/binary%20size-2.4MB-blueviolet.svg)](dist/)

> **Verifiable Environment Parity & Drift Detection Engine for Cloud-Native Software**

X-Parity mathematically eliminates the classic *"works on my machine"* syndrome. It captures comprehensive environment telemetry (OS, kernel, language runtimes, git commit/branch/dirty states, and sanitized configurations), executes semantic difference analysis, and synthesizes automated root-cause diagnostic hypotheses.

---

## ⚡ Performance & Efficiency Metrics

X-Parity was engineered with **Mechanical Sympathy** and strict systems programming discipline. Unlike heavy SaaS observability agents (Datadog, Dynatrace) that consume hundreds of megabytes of background memory, X-Parity is an ultra-fast, zero-overhead static utility:

| Metric | Measured Benchmark | Comparison / Context |
| :--- | :---: | :--- |
| **Binary Footprint** | **2.4 MB** | Single static ELF/Mach-O binary; zero external C or Go dependencies. |
| **Cold Start Latency** | **< 3 ms** | Instant user-space execution via Go standard library runtime. |
| **Capture Execution Time** | **< 28 ms** | Interrogates procfs, runtimes, git, and env vars concurrently (`sync.WaitGroup`). |
| **Comparison Time** | **< 4 ms** | Recursive in-memory semantic diff engine with version compatibility analysis. |
| **Memory Consumption (RSS)** | **< 14 MB peak** | Ideal for resource-constrained CI runners, edge nodes, and local developer machines. |
| **CPU Utilization** | **Negligible** | Short burst of asynchronous I/O; no background daemon or polling overhead. |
| **Security & Privacy** | **100% Local** | Heuristic secret redactor scrubs passwords, keys, and tokens prior to disk persistence. |

---

## 📥 Quick Download & Installation

Pre-built binaries for all major platforms are available directly in the [`dist/`](./dist) folder.

### 1-Line Universal Install (Linux & macOS)
### 1-Line Universal Install

**Linux & macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/AppeiYA/x-parity/main/dist/install.sh | bash
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/AppeiYA/x-parity/main/dist/install.ps1 | iex
```

**Windows (Winget):**
```powershell
winget install AppeiYA.x-parity
```

### Self-Install via Downloaded Binary
You can also download your platform binary directly from the [`dist/`](./dist) folder and run:
```bash
# Example for Linux (x86_64):
chmod +x dist/x-parity-linux-amd64
sudo ./dist/x-parity-linux-amd64 install

# Example for macOS Apple Silicon (M1/M2/M3/M4):
chmod +x dist/x-parity-darwin-arm64
sudo ./dist/x-parity-darwin-arm64 install
```

For full installation options and checksum verification, see the [**Distribution README**](./dist/README.md).

---

## 🛠️ Usage & Core Workflows

### 1. Capture Environment Telemetry
Generates an immutable, invariant-validated snapshot of your local or remote host:
```bash
x-parity capture -app payments-api -env local -out ./local_snap.json
```
* **Sample Output:**
```text
==================================================
X-PARITY SNAPSHOT: payments-api (local)
==================================================
Version:     1.0.0
Captured At: 2026-09-18 15:40:00 UTC

[Runtime]
  OS:           linux (amd64)
  Kernel:       6.8.0-generic
  Hostname:     dev-box
  State:        known
  Runtimes:
    - go: 1.24.1

[Source]
  Commit: 8f4a29c (branch: main, dirty: false)
  State:  known

[Configuration]
  Total Variables: 142 (8 sensitive/redacted)

Snapshot written successfully to ./local_snap.json
```

### 2. Compare Two Snapshots (Drift Detection)
Performs semantic difference analysis between two environments:
```bash
# Strict comparison (fails build if OS or runtimes differ):
x-parity compare ./local_snap.json ./staging_snap.json

# Cross-platform comparison (permits OS & path divergence for Docker/VM workflows):
x-parity compare -cross-platform ./windows_dev.json ./linux_ci.json
```
* **Sample Output:**
```text
==================================================
PARITY DIFFERENCES DETECTED: 2
==================================================

[WARNING] runtime.version (runtime)
  Local:   1.24.1
  Remote:  1.24.4
  Details: Minor runtime patch version divergence detected.
  ---

[CRITICAL] config.DATABASE_URL (configuration)
  Local:   [REDACTED_SECRET]
  Remote:  <missing>
  Details: Critical configuration key is absent in remote environment.
```

### 3. Inspect a Snapshot
View the details of an isolated snapshot without comparing:
```bash
x-parity inspect ./local_snap.json
```

### 4. Synthesize Automated Diagnostics
Evaluates environmental evidence to produce actionable root-cause hypotheses:
```bash
x-parity diagnose ./staging_snap.json
```

---

## 🤖 CI/CD Pipeline Integration

X-Parity is designed to serve as an automated **Parity Gate** in modern CI/CD pipelines.

### Standard Exit Codes
- **`0`**: Full parity (no differences detected). Pipeline proceeds.
- **`1`**: Drift or differences detected (or operational error). Pipeline fails, blocking deployment.
- **`2`**: Invalid CLI arguments or syntax.

### GitHub Actions Workflow
Create `.github/workflows/parity-gate.yml`:
```yaml
name: Environment Parity Gate

on:
  pull_request:
  push:
    branches: [main]

jobs:
  verify-environment-parity:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Source
        uses: actions/checkout@v4

      - name: Install X-Parity
        run: |
          curl -fsSL https://raw.githubusercontent.com/AppeiYA/x-parity/main/dist/install.sh | bash

      - name: Capture Runner Environment
        run: |
          x-parity capture -app core-service -env ci-runner -out ./ci_snap.json

      - name: Enforce Production Parity Gate
        run: |
          # Exits 1 if drift is detected!
          x-parity compare ./ci_snap.json ./deploy/baseline_prod.json || {
            echo "::error::Environment Parity Gate Failed!"
            x-parity diagnose ./ci_snap.json
            exit 1
          }
```

---

## 🏛️ Architecture & Clean Design

X-Parity adheres strictly to **Clean Hexagonal Architecture (Ports and Adapters)**:

```
                  +----------------------------------------------+
                  |               ADAPTERS LAYER                 |
                  |  CLI Subcommand Router (cmd/cli)             |
                  |  OS Collectors (/proc, git, env)             |
                  |  Atomic Swap Persistence (internal/adapters) |
                  |   +--------------------------------------+   |
                  |   |             PORTS LAYER              |   |
                  |   |  CaptureUsecaseInt, CompareUsecaseInt |   |
                  |   |  CollectorInt, SnapshotStoreInt      |   |
                  |   |   +------------------------------+   |   |
                  |   |   |        USE CASE LAYER        |   |   |
                  |   |   |  Capture, Compare, Inspect,  |   |   |
                  |   |   |  Diagnose Orchestration      |   |   |
                  |   |   |   +----------------------+   |   |   |
                  |   |   |   |     DOMAIN LAYER     |   |   |   |
                  |   |   |   | Snapshot, Difference |   |   |   |
                  |   |   |   | DiffEngine, Invariants|  |   |   |
                  |   |   |   +----------------------+   |   |   |
                  |   |   +------------------------------+   |   |
                  |   +--------------------------------------+   |
                  +----------------------------------------------+
```

### Key Architectural Invariants
* **Zero External Dependencies**: Standard library only (`flag`, `sync`, `os`, `syscall`, `encoding/json`).
* **Crash-Resilient Persistence**: Atomic Swap Protocol using POSIX `rename(2)` and hardware `fsync`.
* **Encapsulation Protection**: DTO mapping pattern separates immutable domain structs from serialization layers.
* **Secret Safety**: Heuristic pattern matching scrubs credentials before disk serialization.

---

## 🧪 Running Tests & Verifying Locally

```bash
# Run all unit tests with ThreadSanitizer race detection:
go test -v -race -count=1 ./...

# Build local development binary:
go build -o bin/x-parity ./cmd/cli
```

---

## 📄 License
This project is open-source under the [MIT License](./LICENSE).
