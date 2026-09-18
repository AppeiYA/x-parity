package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/AppeiYA/x-parity/internal/adapters/out/collector"
	"github.com/AppeiYA/x-parity/internal/adapters/out/persistence"
	"github.com/AppeiYA/x-parity/internal/adapters/out/presenter"
	"github.com/AppeiYA/x-parity/internal/domain"
	portin "github.com/AppeiYA/x-parity/internal/ports/in"
	portout "github.com/AppeiYA/x-parity/internal/ports/out"
	"github.com/AppeiYA/x-parity/internal/usecase"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// 1. Initialize Adapters
	termPresenter := presenter.NewTerminalPresenter()
	snapshotStore := persistence.NewFileSnapshotStore()

	configCol := collector.NewConfigCollector()
	systemCol := collector.NewSystemCollector()
	runtimeCol := collector.NewRuntimeCollector()
	collectors := []portout.CollectorInt{configCol, systemCol, runtimeCol}

	diffEngine := domain.NewDiffEngine()

	// 2. Initialize Use Cases
	captureUsecase := usecase.NewCaptureUsecase(collectors, snapshotStore)
	compareUsecase := usecase.NewCompareUsecase(snapshotStore, diffEngine)
	inspectUsecase := usecase.NewInspectUsecase(snapshotStore)
	diagnoseUsecase := usecase.NewDiagnoseUsecase(snapshotStore, diffEngine)

	// 3. Subcommand Routing
	command := os.Args[1]

	switch command {
	case "capture":
		runCapture(ctx, os.Args[2:], captureUsecase, termPresenter)
	case "compare":
		runCompare(ctx, os.Args[2:], compareUsecase, termPresenter)
	case "inspect":
		runInspect(ctx, os.Args[2:], inspectUsecase, termPresenter)
	case "diagnose":
		runDiagnose(ctx, os.Args[2:], diagnoseUsecase, termPresenter)
	case "help", "-h", "--help":
		printUsage()
	default:
		termPresenter.RenderError(fmt.Errorf("unknown command %q", command))
		printUsage()
		os.Exit(1)
	}
}

func runCapture(ctx context.Context, args []string, uc *usecase.CaptureUsecase, p portout.PresenterInt) {
	fs := flag.NewFlagSet("capture", flag.ExitOnError)
	appName := fs.String("app", "", "Application name (required)")
	envName := fs.String("env", "", "Environment name (required)")
	outputPath := fs.String("out", "", "Output snapshot file path (optional)")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: x-parity capture -app <name> -env <env> [-out <file.json>]\n\nFlags:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		p.RenderError(err)
		os.Exit(1)
	}

	if *appName == "" || *envName == "" {
		p.RenderError(fmt.Errorf("both -app and -env flags are required"))
		fs.Usage()
		os.Exit(1)
	}

	cmd := portin.CaptureCommand{
		AppName:     *appName,
		Environment: *envName,
		OutputPath:  *outputPath,
	}

	snapshot, err := uc.Execute(ctx, cmd)
	if err != nil {
		p.RenderError(fmt.Errorf("capture failed: %w", err))
		os.Exit(1)
	}

	if err := p.RenderSnapshot(snapshot); err != nil {
		p.RenderError(err)
	}

	if *outputPath != "" {
		fmt.Printf("Snapshot written successfully to %s\n", *outputPath)
	}
}

func runCompare(ctx context.Context, args []string, uc *usecase.CompareUsecase, p portout.PresenterInt) {
	if len(args) < 2 {
		p.RenderError(fmt.Errorf("compare requires two snapshot file paths"))
		fmt.Fprintf(os.Stderr, "Usage: x-parity compare <local_snapshot.json> <remote_snapshot.json>\n")
		os.Exit(1)
	}

	localPath := args[0]
	remotePath := args[1]

	diffs, err := uc.Execute(ctx, localPath, remotePath)
	if err != nil {
		p.RenderError(fmt.Errorf("comparison failed: %w", err))
		os.Exit(1)
	}

	if err := p.RenderDifferences(diffs); err != nil {
		p.RenderError(err)
	}
}

func runInspect(ctx context.Context, args []string, uc *usecase.InspectUsecase, p portout.PresenterInt) {
	if len(args) < 1 {
		p.RenderError(fmt.Errorf("inspect requires a snapshot file path"))
		fmt.Fprintf(os.Stderr, "Usage: x-parity inspect <snapshot.json>\n")
		os.Exit(1)
	}

	snapshotPath := args[0]

	snapshot, err := uc.Execute(ctx, snapshotPath)
	if err != nil {
		p.RenderError(fmt.Errorf("inspect failed: %w", err))
		os.Exit(1)
	}

	if err := p.RenderSnapshot(snapshot); err != nil {
		p.RenderError(err)
	}
}

func runDiagnose(ctx context.Context, args []string, uc *usecase.DiagnoseUsecase, p portout.PresenterInt) {
	if len(args) < 1 {
		p.RenderError(fmt.Errorf("diagnose requires a snapshot file path"))
		fmt.Fprintf(os.Stderr, "Usage: x-parity diagnose <snapshot.json>\n")
		os.Exit(1)
	}

	snapshotPath := args[0]

	report, err := uc.Execute(ctx, snapshotPath)
	if err != nil {
		p.RenderError(fmt.Errorf("diagnosis failed: %w", err))
		os.Exit(1)
	}

	if err := p.RenderDiagnosis(report); err != nil {
		p.RenderError(err)
	}
}

func printUsage() {
	fmt.Printf(`X-Parity: Environment Drift Detection & Parity Verification CLI

Usage:
  x-parity <command> [arguments]

Commands:
  capture   Capture environment telemetry and generate a snapshot
  compare   Compare two snapshots and report parity differences
  inspect   Inspect and print details of a single snapshot
  diagnose  Analyze a snapshot and produce root-cause hypotheses
  help      Display help information

Examples:
  x-parity capture -app demo-app -env local -out ./local_snap.json
  x-parity compare ./local_snap.json ./remote_snap.json
  x-parity inspect ./local_snap.json
  x-parity diagnose ./local_snap.json
`)
}
