package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/AppeiYA/x-parity/internal/domain"
	portin "github.com/AppeiYA/x-parity/internal/ports/in"
	portout "github.com/AppeiYA/x-parity/internal/ports/out"
)

var _ portin.InspectInt = (*InspectUsecase)(nil)

var (
	ErrEmptySnapshotPath = errors.New("snapshot path cannot be empty")
)

type InspectUsecase struct {
	store portout.SnapshotStore
}

func NewInspectUsecase(store portout.SnapshotStore) *InspectUsecase {
	return &InspectUsecase{
		store: store,
	}
}

func (u *InspectUsecase) Execute(ctx context.Context, snapshotPath string) (*domain.Snapshot, error) {
	if snapshotPath == "" {
		return nil, ErrEmptySnapshotPath
	}
	if u.store == nil {
		return nil, errors.New("snapshot store is not configured")
	}

	snap, err := u.store.Load(ctx, snapshotPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load snapshot from %q: %w", snapshotPath, err)
	}
	if snap == nil {
		return nil, fmt.Errorf("loaded snapshot from %q is nil", snapshotPath)
	}

	return snap, nil
}

