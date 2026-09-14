package portout

import "github.com/AppeiYA/x-parity/internal/domain"

type PresenterInt interface {
	RenderSnapshot(snapshot *domain.Snapshot) error
	RenderDifferences(diffs []domain.Difference) error
	RenderDiagnosis(report *domain.DiagnosisReport) error
	RenderError(err error) error
}