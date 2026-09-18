# X-Parity Binary Distribution & Installation Guide

Pre-compiled, zero-dependency standalone binaries for **X-Parity**: the high-performance environment drift detection and parity verification engine.

---

## 📦 Pre-Compiled Release Binaries

All binaries are compiled statically with `CGO_ENABLED=0` and stripped of DWARF symbols (`-ldflags="-s -w"`), resulting in microscopic binary footprints (~2.4MB) with zero runtime dependencies.

| Platform | Architecture | Binary File | Description |
| :--- | :--- | :--- | :--- |
| **Linux** | `amd64` (x86_64) | [`x-parity-linux-amd64`](./x-parity-linux-amd64) | Standard 64-bit Linux (Ubuntu, Debian, Fedora, Arch, Alpine, RHEL) |
| **Linux** | `arm64` (aarch64) | [`x-parity-linux-arm64`](./x-parity-linux-arm64) | 64-bit ARM Linux (AWS Graviton, Raspberry Pi 4/5, Ampere) |
| **macOS** | `arm64` | [`x-parity-darwin-arm64`](./x-parity-darwin-arm64) | Apple Silicon Macs (M1, M2, M3, M4) |
| **macOS** | `amd64` (x86_64) | [`x-parity-darwin-amd64`](./x-parity-darwin-amd64) | Intel-based Macs |
| **Windows** | `amd64` (x86_64) | [`x-parity-windows-amd64.exe`](./x-parity-windows-amd64.exe) | Standard 64-bit Windows 10/11 / Windows Server |

---

## 🚀 Installation Methods

### Method 1: Self-Install (Zero Configuration)
Download the binary for your platform, make it executable, and run the built-in `install` command. It will automatically detect permissions and place `x-parity` into `/usr/local/bin` (or `~/.local/bin`):

```bash
# Example for Linux (x86_64):
chmod +x x-parity-linux-amd64
sudo ./x-parity-linux-amd64 install

# Example for macOS (Apple Silicon M1/M2/M3/M4):
chmod +x x-parity-darwin-arm64
sudo ./x-parity-darwin-arm64 install
```

> **Note**: If you do not have `sudo` privileges, run `./x-parity-<os>-<arch> install` without `sudo`. The binary will automatically install to `~/.local/bin/x-parity`.

---

### Method 2: Universal Installer Script (`install.sh`)
You can use the included `install.sh` script, which automatically detects your operating system and CPU architecture:

```bash
# Run locally from this directory:
chmod +x install.sh
./install.sh

# Or install via one-liner from GitHub:
curl -fsSL https://raw.githubusercontent.com/AppeiYA/x-parity/main/dist/install.sh | bash
```

---

### Method 3: Windows PowerShell Installation
### Method 3: Windows Automated Installation

#### Option A: PowerShell 1-Liner (Recommended)
Open PowerShell and run:
```powershell
irm https://raw.githubusercontent.com/AppeiYA/x-parity/main/dist/install.ps1 | iex
```

#### Option B: Windows Package Manager (Winget)
```powershell
# 1. Download binary (or use the one in this folder)
# 2. Place in your user programs folder
winget install AppeiYA.x-parity
```

#### Option C: Chocolatey
```powershell
choco install x-parity
```

#### Option D: Manual Copy
```powershell
mkdir -Force "$env:LOCALAPPDATA\Programs\x-parity"
Copy-Item x-parity-windows-amd64.exe "$env:LOCALAPPDATA\Programs\x-parity\x-parity.exe"

# 3. Add to user PATH (if not already added)
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*x-parity*") {
    [Environment]::SetEnvironmentVariable("Path", "$userPath;$env:LOCALAPPDATA\Programs\x-parity", "User")
}
```

---

## 🔒 Cryptographic Verification (SHA-256)

Verify the integrity of downloaded binaries against `checksums.sha256`:

```bash
sha256sum -c checksums.sha256
```

---

## ⚡ Quick Start: Using X-Parity

Once installed, `x-parity` is available globally across all terminals.

### 1. Capture Environment Telemetry
Interrogates host OS, kernel, language runtimes, git commit/branch/dirty status, and environment variables. Sensitive secrets are **automatically redacted**.

```bash
x-parity capture -app payment-service -env local -out ./local_snap.json
```

### 2. Inspect a Snapshot
View captured metadata, runtime details, and scrubbed configuration counts:

```bash
x-parity inspect ./local_snap.json
```

### 3. Compare Two Environments (Detect Drift)
Compare two snapshots (e.g., local developer laptop vs. staging or production):

```bash
# Standard strict comparison (fails if OS, runtimes, or configs diverge):
x-parity compare ./local_snap.json ./staging_snap.json

# Cross-platform comparison (for Windows/macOS laptop vs. Docker container or Linux CI):
# Permits OS and path format divergence (INFO), while strictly enforcing language runtimes and configs
x-parity compare -cross-platform ./windows_dev.json ./linux_ci.json
```
* **Exit code `0`**: Environments are fully aligned.
* **Exit code `1`**: Parity divergence or drift detected.
* **Exit code `0`**: Environments are fully aligned (or only have permitted cross-platform variances).
* **Exit code `1`**: Parity divergence or breaking drift detected (fails CI/CD build).

### 4. Synthesize Automated Diagnostics
Analyze environmental evidence to produce prioritized root-cause hypotheses:

```bash
x-parity diagnose ./staging_snap.json
```

---

## 🤖 CI/CD Integration & Automation

X-Parity is purpose-built for CI/CD pipelines. It is lightweight, fast, and uses standard Unix exit codes.

### Exit Codes Reference
| Code | Meaning | Pipeline Behavior |
| :---: | :--- | :--- |
| `0` | Clean match / snapshots in complete parity | Step passes (CI continues) |
| `1` | Parity differences detected OR runtime error | Step fails (Blocks PR / deployment) |
| `2` | Flag or argument syntax error | Step fails |

### GitHub Actions Workflow Example

Add this step to `.github/workflows/deploy.yml` to prevent unaligned code or drifted staging environments from reaching production:

```yaml
name: Environment Parity Gate

on:
  pull_request:
  push:
    branches: [main]

jobs:
  verify-parity:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v4

      - name: Install X-Parity
        run: |
          curl -fsSL https://raw.githubusercontent.com/AppeiYA/x-parity/main/dist/install.sh | bash

      - name: Capture Runner Environment
        run: |
          x-parity capture -app backend-api -env ci-runner -out ./ci_snap.json

      - name: Verify Parity against Production Baseline
        run: |
          # Fails the build if drift is detected!
          x-parity compare ./ci_snap.json ./deploy/baseline_prod.json || {
            echo "::error::Parity verification failed between CI and Production!"
            x-parity diagnose ./ci_snap.json
            exit 1
          }
```

### GitLab CI Example (`.gitlab-ci.yml`)

```yaml
parity_check:
  stage: test
  image: alpine:latest
  before_script:
    - apk add --no-cache bash curl
    - curl -fsSL https://raw.githubusercontent.com/AppeiYA/x-parity/main/dist/install.sh | bash
  script:
    - x-parity capture -app core-service -env gitlab-ci -out ./ci_snap.json
    - x-parity compare ./ci_snap.json ./config/prod_baseline.json
```
