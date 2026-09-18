package collector

import (
	"context"
	"os/exec"
	"runtime"
	"strings"

	"github.com/AppeiYA/x-parity/internal/domain"
	portout "github.com/AppeiYA/x-parity/internal/ports/out"
)

var _ portout.CollectorInt = (*RuntimeCollector)(nil)

type RuntimeCollector struct{}

func NewRuntimeCollector() *RuntimeCollector {
	return &RuntimeCollector{}
}

func (rc *RuntimeCollector) Name() string {
	return portout.CollectorRuntime
}

func (rc *RuntimeCollector) Collect(ctx context.Context, s *domain.Snapshot) error {
	rt := s.Runtime()
	if rt == nil {
		var err error
		rt, err = domain.NewRuntime(runtime.GOOS, runtime.GOARCH, "", "", make(map[string]string), domain.StateKnown)
		if err != nil {
			return err
		}
		s.SetRuntime(rt)
	}
	if rt.Runtimes == nil {
		rt.Runtimes = make(map[string]string)
	}

	// 1. Go Runtime
	if _, err := exec.LookPath("go"); err == nil {
		if out, err := exec.CommandContext(ctx, "go", "version").Output(); err == nil {
			parts := strings.Fields(string(out))
			if len(parts) >= 3 {
				rt.Runtimes["go"] = strings.TrimPrefix(parts[2], "go")
			}
		}
	}

	// 2. Node.js Runtime
	if _, err := exec.LookPath("node"); err == nil {
		if out, err := exec.CommandContext(ctx, "node", "--version").Output(); err == nil {
			rt.Runtimes["node"] = strings.TrimSpace(strings.TrimPrefix(string(out), "v"))
		}
	}

	// 3. Python Runtime (python3 on Unix, python or py on Windows)
	pythonBins := []string{"python3", "python", "py"}
	for _, bin := range pythonBins {
		if _, err := exec.LookPath(bin); err == nil {
			if out, err := exec.CommandContext(ctx, bin, "--version").Output(); err == nil {
				parts := strings.Fields(string(out))
				if len(parts) >= 2 {
					rt.Runtimes["python"] = parts[1]
					break
				}
			}
		}
	}

	// 4. .NET Runtime
	if _, err := exec.LookPath("dotnet"); err == nil {
		if out, err := exec.CommandContext(ctx, "dotnet", "--version").Output(); err == nil {
			rt.Runtimes["dotnet"] = strings.TrimSpace(string(out))
		}
	}

	// 5. PowerShell
	pwshBins := []string{"pwsh", "powershell"}
	for _, bin := range pwshBins {
		if _, err := exec.LookPath(bin); err == nil {
			if out, err := exec.CommandContext(ctx, bin, "-NoProfile", "-Command", "$PSVersionTable.PSVersion.ToString()").Output(); err == nil {
				ver := strings.TrimSpace(string(out))
				if ver != "" {
					rt.Runtimes["powershell"] = ver
					break
				}
			}
		}
	}

	return nil
}