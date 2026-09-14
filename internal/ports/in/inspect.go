package portin

import (
	"context"

	"github.com/AppeiYA/x-parity/internal/domain"
)

type InspectInt interface {
	Execute(ctx context.Context, snapshotPath string) (*domain.Snapshot, error)
}