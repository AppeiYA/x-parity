package presenter_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/AppeiYA/x-parity/internal/adapters/out/presenter"
	"github.com/AppeiYA/x-parity/internal/domain"
)

func TestTerminalPresenter_RenderSnapshot(t *testing.T) {
	var out, errBuf bytes.Buffer
	p := presenter.NewCustomTerminalPresenter(&out, &errBuf)

	app, _ := domain.NewApplication("app-test", "git@github.com:org/repo.git")
	env, _ := domain.NewEnvironment("dev", "local")
	snap, _ := domain.NewSnapshot("1.0.0", time.Now().UTC(), app, env)

	rt, _ := domain.NewRuntime("linux", "arm64", "6.1", "host-1", map[string]string{"go": "1.25"}, domain.StateKnown)
	snap.SetRuntime(rt)

	if err := p.RenderSnapshot(snap); err != nil {
		t.Fatalf("RenderSnapshot failed: %v", err)
	}

	result := out.String()
	if !strings.Contains(result, "X-PARITY SNAPSHOT: app-test (dev)") {
		t.Errorf("expected banner in output, got: %s", result)
	}
	if !strings.Contains(result, "linux (arm64)") {
		t.Errorf("expected runtime info, got: %s", result)
	}
	if !strings.Contains(result, "go: 1.25") {
		t.Errorf("expected go runtime info, got: %s", result)
	}
}

func TestTerminalPresenter_RenderDifferences(t *testing.T) {
	var out, errBuf bytes.Buffer
	p := presenter.NewCustomTerminalPresenter(&out, &errBuf)

	diff, _ := domain.NewDifference(
		domain.CatRuntime,
		"runtime.os",
		"linux",
		"darwin",
		domain.SeverityCritical,
		"OS mismatch",
	)

	if err := p.RenderDifferences([]domain.Difference{*diff}); err != nil {
		t.Fatalf("RenderDifferences failed: %v", err)
	}

	result := out.String()
	if !strings.Contains(result, "[CRITICAL] runtime.os") {
		t.Errorf("expected critical diff badge, got: %s", result)
	}
}

func TestTerminalPresenter_RenderDiagnosis(t *testing.T) {
	var out, errBuf bytes.Buffer
	p := presenter.NewCustomTerminalPresenter(&out, &errBuf)

	report := &domain.DiagnosisReport{
		TargetSnapshot: "/path/snap.json",
		GeneratedAt:    time.Now().UTC(),
		Summary:        "Environment disparity detected",
		Hypotheses: []domain.Hypothesis{
			{
				Title:       "Missing Go Runtime",
				LikelyCause: "Go is not installed on target host",
				Confidence:  domain.TierDeterministic,
				Evidence:    []string{"runtime.runtimes.go missing"},
				Remediation: "Install go 1.25 on remote environment",
			},
		},
	}

	if err := p.RenderDiagnosis(report); err != nil {
		t.Fatalf("RenderDiagnosis failed: %v", err)
	}

	result := out.String()
	if !strings.Contains(result, "X-PARITY DIAGNOSIS REPORT") {
		t.Errorf("expected diagnosis header, got: %s", result)
	}
	if !strings.Contains(result, "[DETERMINISTIC CONFIDENCE] Missing Go Runtime") {
		t.Errorf("expected hypothesis title, got: %s", result)
	}
}

func TestTerminalPresenter_RenderError(t *testing.T) {
	var out, errBuf bytes.Buffer
	p := presenter.NewCustomTerminalPresenter(&out, &errBuf)

	testErr := errors.New("something went wrong")
	_ = p.RenderError(testErr)

	if !strings.Contains(errBuf.String(), "Error: something went wrong") {
		t.Errorf("expected error output in errBuf, got: %s", errBuf.String())
	}
}
