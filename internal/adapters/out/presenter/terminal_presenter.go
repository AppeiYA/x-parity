package presenter

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/AppeiYA/x-parity/internal/domain"
	portout "github.com/AppeiYA/x-parity/internal/ports/out"
)

var _ portout.PresenterInt = (*TerminalPresenter)(nil)

type TerminalPresenter struct {
	out io.Writer
	err io.Writer
}

func NewTerminalPresenter() *TerminalPresenter {
	return &TerminalPresenter{
		out: os.Stdout,
		err: os.Stderr,
	}
}

func NewCustomTerminalPresenter(out, err io.Writer) *TerminalPresenter {
	return &TerminalPresenter{
		out: out,
		err: err,
	}
}

func (tp *TerminalPresenter) RenderSnapshot(snap *domain.Snapshot) error {
	if snap == nil {
		return fmt.Errorf("cannot render nil snapshot")
	}

	fmt.Fprintln(tp.out, "==================================================")
	appName := "Unknown"
	envName := "Unknown"
	if snap.Application() != nil {
		appName = snap.Application().Name
	}
	if snap.Environment() != nil {
		envName = snap.Environment().Name
	}
	fmt.Fprintf(tp.out, "X-PARITY SNAPSHOT: %s (%s)\n", appName, envName)
	fmt.Fprintln(tp.out, "==================================================")
	fmt.Fprintf(tp.out, "Version:     %s\n", snap.Version())
	fmt.Fprintf(tp.out, "Captured At: %s\n", snap.CapturedAt().Format("2006-01-02 15:04:05 UTC"))
	if snap.Application() != nil && snap.Application().Repository != "" {
		fmt.Fprintf(tp.out, "Repository:  %s\n", snap.Application().Repository)
	}

	if rt := snap.Runtime(); rt != nil {
		fmt.Fprintln(tp.out, "\n[Runtime]")
		fmt.Fprintf(tp.out, "  OS:           %s (%s)\n", rt.OS, rt.Architecture)
		if rt.Kernel != "" {
			fmt.Fprintf(tp.out, "  Kernel:       %s\n", rt.Kernel)
		}
		if rt.Hostname != "" {
			fmt.Fprintf(tp.out, "  Hostname:     %s\n", rt.Hostname)
		}
		fmt.Fprintf(tp.out, "  State:        %s\n", rt.State)
		if len(rt.Runtimes) > 0 {
			fmt.Fprintln(tp.out, "  Runtimes:")
			for lang, ver := range rt.Runtimes {
				fmt.Fprintf(tp.out, "    - %s: %s\n", lang, ver)
			}
		}
	}

	if src := snap.Source(); src != nil {
		fmt.Fprintln(tp.out, "\n[Source]")
		fmt.Fprintf(tp.out, "  Commit: %s (branch: %s, dirty: %t)\n", src.Commit, src.Branch, src.Dirty)
		fmt.Fprintf(tp.out, "  State:  %s\n", src.State)
	}

	if cfg := snap.Configuration(); cfg != nil {
		fmt.Fprintln(tp.out, "\n[Configuration]")
		sensitiveCount := 0
		for _, v := range cfg.Variables {
			if v != nil && v.Sensitive {
				sensitiveCount++
			}
		}
		fmt.Fprintf(tp.out, "  Total Variables: %d (%d sensitive/redacted)\n", len(cfg.Variables), sensitiveCount)
	}

	if deps := snap.Dependencies(); len(deps) > 0 {
		fmt.Fprintln(tp.out, "\n[Dependencies]")
		for _, d := range deps {
			fmt.Fprintf(tp.out, "  - %s (%s)", d.Name, d.Type)
			if d.Version != "" {
				fmt.Fprintf(tp.out, " version: %s", d.Version)
			}
			if d.Host != "" {
				fmt.Fprintf(tp.out, " host: %s:%s", d.Host, d.Port)
			}
			fmt.Fprintln(tp.out)
		}
	}

	fmt.Fprintln(tp.out)
	return nil
}

func (tp *TerminalPresenter) RenderDifferences(diffs []domain.Difference) error {
	if len(diffs) == 0 {
		fmt.Fprintln(tp.out, "✓ No parity differences detected. Environments are fully aligned.")
		return nil
	}

	fmt.Fprintln(tp.out, "==================================================")
	fmt.Fprintf(tp.out, "PARITY DIFFERENCES DETECTED: %d\n", len(diffs))
	fmt.Fprintln(tp.out, "==================================================")

	for i, d := range diffs {
		sevBadge := strings.ToUpper(string(d.Severity()))
		fmt.Fprintf(tp.out, "\n[%s] %s (%s)\n", sevBadge, d.Path(), d.Category())
		fmt.Fprintf(tp.out, "  Local:   %v\n", d.Local())
		fmt.Fprintf(tp.out, "  Remote:  %v\n", d.Remote())
		if d.Message() != "" {
			fmt.Fprintf(tp.out, "  Details: %s\n", d.Message())
		}
		if i < len(diffs)-1 {
			fmt.Fprintln(tp.out, "  ---")
		}
	}
	fmt.Fprintln(tp.out)
	return nil
}

func (tp *TerminalPresenter) RenderDiagnosis(report *domain.DiagnosisReport) error {
	if report == nil {
		return fmt.Errorf("cannot render nil diagnosis report")
	}

	fmt.Fprintln(tp.out, "==================================================")
	fmt.Fprintln(tp.out, "X-PARITY DIAGNOSIS REPORT")
	fmt.Fprintln(tp.out, "==================================================")
	if report.TargetSnapshot != "" {
		fmt.Fprintf(tp.out, "Target:       %s\n", report.TargetSnapshot)
	}
	fmt.Fprintf(tp.out, "Generated At: %s\n", report.GeneratedAt.Format("2006-01-02 15:04:05 UTC"))
	if report.Summary != "" {
		fmt.Fprintf(tp.out, "Summary:      %s\n", report.Summary)
	}

	if len(report.Differences) > 0 {
		fmt.Fprintf(tp.out, "\nRoot Differences (%d):\n", len(report.Differences))
		for _, d := range report.Differences {
			fmt.Fprintf(tp.out, "  - [%s] %s: local=%v, remote=%v\n", d.Severity(), d.Path(), d.Local(), d.Remote())
		}
	}

	if len(report.Hypotheses) > 0 {
		fmt.Fprintln(tp.out, "\nHypotheses & Remediations:")
		for i, h := range report.Hypotheses {
			fmt.Fprintf(tp.out, "\n%d. [%s CONFIDENCE] %s\n", i+1, strings.ToUpper(h.Confidence.String()), h.Title)
			if h.LikelyCause != "" {
				fmt.Fprintf(tp.out, "   Likely Cause: %s\n", h.LikelyCause)
			}
			if len(h.Evidence) > 0 {
				fmt.Fprintln(tp.out, "   Evidence:")
				for _, ev := range h.Evidence {
					fmt.Fprintf(tp.out, "     • %s\n", ev)
				}
			}
			if h.Remediation != "" {
				fmt.Fprintf(tp.out, "   Remediation:  %s\n", h.Remediation)
			}
		}
	}
	fmt.Fprintln(tp.out)
	return nil
}

func (tp *TerminalPresenter) RenderError(err error) error {
	if err == nil {
		return nil
	}
	fmt.Fprintf(tp.err, "Error: %v\n", err)
	return nil
}
