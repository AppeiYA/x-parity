//go:build !linux && !windows && !darwin

package collector

import (
	"context"
	"os/exec"
	"strings"

	"github.com/AppeiYA/x-parity/internal/domain"
)

func collectPlatformTelemetry(ctx context.Context, rt *domain.Runtime) error {
	if out, err := exec.CommandContext(ctx, "uname", "-r").Output(); err == nil {
		rt.Kernel = strings.TrimSpace(string(out))
	} else {
		rt.Kernel = "generic"
	}
	return nil
}

