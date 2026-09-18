package collector

import (
	"context"
	"os"
	"runtime"

	"github.com/AppeiYA/x-parity/internal/domain"
	portout "github.com/AppeiYA/x-parity/internal/ports/out"
)

var _ portout.CollectorInt = (*SystemCollector)(nil)

type SystemCollector struct{}

func NewSystemCollector() *SystemCollector {
	return &SystemCollector{}
}

func (sc *SystemCollector) Name() string {
	return portout.CollectorSystem
}

func (sc *SystemCollector) Collect(ctx context.Context, s *domain.Snapshot) error {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	rt := s.Runtime()
	if rt == nil {
		var err error
		rt, err = domain.NewRuntime(runtime.GOOS, runtime.GOARCH, "", hostname, make(map[string]string), domain.StateKnown)
		if err != nil {
			return err
		}
		s.SetRuntime(rt)
	} else {
		rt.OS = runtime.GOOS
		rt.Architecture = runtime.GOARCH
		rt.Hostname = hostname
		rt.State = domain.StateKnown
		if rt.Runtimes == nil {
			rt.Runtimes = make(map[string]string)
		}
	}

	return nil
	return collectPlatformTelemetry(ctx, rt)
}