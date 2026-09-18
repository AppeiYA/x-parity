package portout

import (
	"context"

	"github.com/AppeiYA/x-parity/internal/domain"
)

const (
	CollectorSystem        = "system"
	CollectorRuntime       = "runtime"
	CollectorGit           = "git"
	CollectorSource        = "source"
	CollectorEnv           = "env"
	CollectorConfig        = "config"
	CollectorConfiguration = "configuration"
	CollectorBuild         = "build"
	CollectorDeployment    = "deployment"
)

type CollectorInt interface {
	Name() string
	Collect(ctx context.Context, snapshot *domain.Snapshot) error
}