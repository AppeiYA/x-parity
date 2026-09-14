package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/AppeiYA/x-parity/internal/domain"
	portin "github.com/AppeiYA/x-parity/internal/ports/in"
	portout "github.com/AppeiYA/x-parity/internal/ports/out"
)

var _ portin.CompareInt = (*CompareUsecase)(nil)

var (
	ErrEmptyLocalPath  = errors.New("local snapshot path cannot be empty")
	ErrEmptyRemotePath = errors.New("remote snapshot path cannot be empty")
	ErrNilSnapshot     = errors.New("loaded snapshot is nil")
)

type CompareUsecase struct {
	store      portout.SnapshotStore
	diffEngine *domain.DiffEngine
}

func NewCompareUsecase(store portout.SnapshotStore, diffEngine *domain.DiffEngine) *CompareUsecase {
	return &CompareUsecase{
		store:      store,
		diffEngine: diffEngine,
	}
}

func (c *CompareUsecase) Execute(ctx context.Context, localPath, remotePath string) ([]domain.Difference, error) {
	if localPath == "" {
		return nil, ErrEmptyLocalPath
	}
	if remotePath == "" {
		return nil, ErrEmptyRemotePath
	}

	localSnap, err := c.store.Load(ctx, localPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load local snapshot from %q: %w", localPath, err)
	}
	if localSnap == nil {
		return nil, fmt.Errorf("local snapshot (%s): %w", localPath, ErrNilSnapshot)
	}

	remoteSnap, err := c.store.Load(ctx, remotePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load remote snapshot from %q: %w", remotePath, err)
	}
	if remoteSnap == nil {
		return nil, fmt.Errorf("remote snapshot (%s): %w", remotePath, ErrNilSnapshot)
	}

	diffs, err := c.diffEngine.Compare(*localSnap, *remoteSnap)
	if err != nil {
		return nil, fmt.Errorf("diff comparison failed: %w", err)
	}
	return diffs, nil
}