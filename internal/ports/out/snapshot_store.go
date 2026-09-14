package portout

import (
	"context"

	"github.com/AppeiYA/x-parity/internal/domain"
)

type SnapshotStore interface {
	Save(ctx context.Context, path string, snapshot *domain.Snapshot) error
	Load(ctx context.Context, path string) (*domain.Snapshot, error)
}