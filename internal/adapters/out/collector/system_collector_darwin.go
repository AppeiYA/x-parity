//go:build darwin

package collector

import (
	"context"
	"os/exec"
	"strings"

	"github.com/AppeiYA/x-parity/internal/domain"
)

func collectPlatformTelemetry(ctx context.Context, rt *domain.Runtime) error {
	// 1. Interrogate macOS ProductVersion and BuildVersion via sw_vers
	var macVersion, buildVersion string
	if out, err := exec.CommandContext(ctx, "sw_vers", "-productVersion").Output(); err == nil {
		macVersion = strings.TrimSpace(string(out))
	}
	if out, err := exec.CommandContext(ctx, "sw_vers", "-buildVersion").Output(); err == nil {
		buildVersion = strings.TrimSpace(string(out))
	}

	// 2. Interrogate Darwin kernel version via sysctl
	var darwinKernel string
	if out, err := exec.CommandContext(ctx, "sysctl", "-n", "kern.osrelease").Output(); err == nil {
		darwinKernel = strings.TrimSpace(string(out))
	}

	if macVersion != "" {
		kernelStr := "macOS " + macVersion
		if buildVersion != "" {
			kernelStr += " (" + buildVersion + ")"
		}
		if darwinKernel != "" {
			kernelStr += " Darwin " + darwinKernel
		}
		rt.Kernel = kernelStr
	} else if darwinKernel != "" {
		rt.Kernel = "Darwin " + darwinKernel
	}

	// 3. Detect Rosetta 2 translation state
	if out, err := exec.CommandContext(ctx, "sysctl", "-n", "sysctl.proc_translated").Output(); err == nil {
		if strings.TrimSpace(string(out)) == "1" {
			rt.Virtualization = "rosetta2-translation"
		}
	}

	return nil
}

