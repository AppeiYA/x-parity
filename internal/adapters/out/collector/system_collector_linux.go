//go:build linux

package collector

import (
	"bufio"
	"context"
	"os"
	"strings"

	"github.com/AppeiYA/x-parity/internal/domain"
)

func collectPlatformTelemetry(ctx context.Context, rt *domain.Runtime) error {
	// 1. Interrogate Linux Kernel Version
	var kernelVersion string
	if data, err := os.ReadFile("/proc/sys/kernel/osrelease"); err == nil {
		kernelVersion = strings.TrimSpace(string(data))
	} else if data, err := os.ReadFile("/proc/version"); err == nil {
		fields := strings.Fields(string(data))
		if len(fields) >= 3 {
			kernelVersion = fields[2]
		}
	}

	// 2. Interrogate Distribution (/etc/os-release or /usr/lib/os-release)
	distro := parseOSRelease()
	if distro != "" && kernelVersion != "" {
		rt.Kernel = kernelVersion + " (" + distro + ")"
	} else if kernelVersion != "" {
		rt.Kernel = kernelVersion
	} else if distro != "" {
		rt.Kernel = distro
	}

	// 3. Detect WSL2 (Windows Subsystem for Linux)
	// WSL2 is Linux under the hood, so OS remains "linux", but Virtualization aids debugging
	if isWSL2() {
		rt.Virtualization = "wsl2"
	} else if isContainer() {
		rt.Virtualization = "container"
	}

	return nil
}

func parseOSRelease() string {
	files := []string{"/etc/os-release", "/usr/lib/os-release"}
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			continue
		}
		defer f.Close()

		var prettyName, name, version string
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				prettyName = cleanQuotes(strings.TrimPrefix(line, "PRETTY_NAME="))
			} else if strings.HasPrefix(line, "NAME=") {
				name = cleanQuotes(strings.TrimPrefix(line, "NAME="))
			} else if strings.HasPrefix(line, "VERSION_ID=") {
				version = cleanQuotes(strings.TrimPrefix(line, "VERSION_ID="))
			}
		}

		if prettyName != "" {
			return prettyName
		}
		if name != "" {
			if version != "" {
				return name + " " + version
			}
			return name
		}
	}
	return ""
}

func isWSL2() bool {
	if _, err := os.Stat("/proc/sys/fs/binfmt_misc/WSLCommandLine"); err == nil {
		return true
	}
	if data, err := os.ReadFile("/proc/version"); err == nil {
		lower := strings.ToLower(string(data))
		if strings.Contains(lower, "microsoft") || strings.Contains(lower, "wsl") {
			return true
		}
	}
	return false
}

func isContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	if _, err := os.Stat("/run/.containerenv"); err == nil {
		return true
	}
	if data, err := os.ReadFile("/proc/1/cgroup"); err == nil {
		content := string(data)
		if strings.Contains(content, "docker") || strings.Contains(content, "containerd") || strings.Contains(content, "kubepods") || strings.Contains(content, "lxc") {
			return true
		}
	}
	return false
}

func cleanQuotes(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "\"")
	s = strings.TrimSuffix(s, "\"")
	s = strings.TrimPrefix(s, "'")
	s = strings.TrimSuffix(s, "'")
	return strings.TrimSpace(s)
}

