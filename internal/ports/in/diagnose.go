package portin

import (
	"context"

	"github.com/AppeiYA/x-parity/internal/domain"
)

type DiagnoseInt interface {
	Execute(ctx context.Context, snapshotPath string) (*domain.DiagnosisReport, error)
}