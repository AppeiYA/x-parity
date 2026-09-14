package portin

import (
	"context"

	"github.com/AppeiYA/x-parity/internal/domain"
)

type CompareInt interface {
	Execute(ctx context.Context, localPath, remotePath string) ([]domain.Difference, error)
}