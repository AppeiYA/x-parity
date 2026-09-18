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

	if out, err := exec.CommandContext(ctx, "go", "version").Output(); err == nil {
		parts := strings.Fields(string(out))
		if len(parts) >= 3 {
			rt.Runtimes["go"] = strings.TrimPrefix(parts[2], "go")
		}
	}

	if out, err := exec.CommandContext(ctx, "node", "--version").Output(); err == nil {
		rt.Runtimes["node"] = strings.TrimSpace(strings.TrimPrefix(string(out), "v"))
	}

	return nil
}