package portin

import (
	"context"

	"github.com/AppeiYA/x-parity/internal/domain"
)

type CaptureCommand struct {
	AppName string
	Environment string
	OutputPath string
}

type CaptureInt interface {
	Execute(ctx context.Context, cmd CaptureCommand) (*domain.Snapshot, error)
}