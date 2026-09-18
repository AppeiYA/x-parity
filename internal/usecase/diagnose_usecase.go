package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AppeiYA/x-parity/internal/domain"
	portin "github.com/AppeiYA/x-parity/internal/ports/in"
	portout "github.com/AppeiYA/x-parity/internal/ports/out"
)

var _ portin.DiagnoseInt = (*DiagnoseUsecase)(nil)

type DiagnoseUsecase struct {
	store      portout.SnapshotStore
	diffEngine *domain.DiffEngine
}

func NewDiagnoseUsecase(store portout.SnapshotStore, diffEngine *domain.DiffEngine) *DiagnoseUsecase {
	return &DiagnoseUsecase{
		store:      store,
		diffEngine: diffEngine,
	}
}

func (u *DiagnoseUsecase) Execute(ctx context.Context, snapshotPath string) (*domain.DiagnosisReport, error) {
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

	var hypotheses []domain.Hypothesis

	// 1. Check Runtime state
	if snap.Runtime() == nil {
		hypotheses = append(hypotheses, domain.Hypothesis{
			Title:       "Missing Runtime Telemetry",
			LikelyCause: "Runtime collector did not run or failed to collect telemetry",
			Confidence:  domain.TierDeterministic,
			Evidence:    []string{"Runtime subsystem is nil in snapshot"},
			Remediation: "Verify collector registration and rerun capture",
		})
	} else if snap.Runtime().State == domain.StateUnavailable {
		hypotheses = append(hypotheses, domain.Hypothesis{
			Title:       "Degraded Runtime Telemetry",
			LikelyCause: "System or runtime collector encountered errors during capture",
			Confidence:  domain.TierDeterministic,
			Evidence:    []string{"Runtime evidence state is unavailable"},
			Remediation: "Check host permissions and command availability (e.g. go, node, uname)",
		})
	}

	// 2. Check Source state
	if snap.Source() != nil {
		if snap.Source().Dirty {
			hypotheses = append(hypotheses, domain.Hypothesis{
				Title:       "Dirty Working Tree Captured",
				LikelyCause: "Snapshot was generated from uncommitted local modifications",
				Confidence:  domain.TierDeterministic,
				Evidence:    []string{fmt.Sprintf("Git commit %s marked as dirty", snap.Source().Commit)},
				Remediation: "Commit or stash pending changes before capturing a release baseline",
			})
		}
		if snap.Source().State == domain.StateUnavailable {
			hypotheses = append(hypotheses, domain.Hypothesis{
				Title:       "VCS Metadata Missing",
				LikelyCause: "Git repository or CLI is not accessible in environment",
				Confidence:  domain.TierDerived,
				Evidence:    []string{"Source evidence state is unavailable"},
				Remediation: "Execute capture inside a cloned git repository with git installed",
			})
		}
	}

	// 3. Check Configuration state
	if snap.Configuration() != nil {
		var missingVars []string
		for k, v := range snap.Configuration().Variables {
			if v != nil && v.State == domain.StateMissing {
				missingVars = append(missingVars, k)
			}
		}
		if len(missingVars) > 0 {
			hypotheses = append(hypotheses, domain.Hypothesis{
				Title:       "Missing Environment Variables",
				LikelyCause: "Application configuration variables are declared missing",
				Confidence:  domain.TierDerived,
				Evidence:    missingVars,
				Remediation: "Supply the missing configuration variables in the execution environment",
			})
		}
	}

	summary := "Snapshot environment telemetry is healthy."
	if len(hypotheses) > 0 {
		summary = fmt.Sprintf("Identified %d potential issue(s) or evidence degradation(s).", len(hypotheses))
	}

	return &domain.DiagnosisReport{
		TargetSnapshot: snapshotPath,
		GeneratedAt:    time.Now().UTC(),
		Differences:    nil,
		Hypotheses:     hypotheses,
		Summary:        summary,
	}, nil
}

