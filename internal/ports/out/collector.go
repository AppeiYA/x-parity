package portout

import (
	"context"

	"github.com/AppeiYA/x-parity/internal/domain"
)

type CollectorInt interface {
	Name() string
	Collect(ctx context.Context, snapshot *domain.Snapshot) error
}