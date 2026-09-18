package collector_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/AppeiYA/x-parity/internal/adapters/out/collector"
	"github.com/AppeiYA/x-parity/internal/domain"
)

func createTestSnapshot(t *testing.T) *domain.Snapshot {
	t.Helper()
	app, _ := domain.NewApplication("myapp", "")
	env, _ := domain.NewEnvironment("local", "dev")
	snap, err := domain.NewSnapshot("1.0.0", time.Now().UTC(), app, env)
	if err != nil {
		t.Fatalf("failed to create snapshot: %v", err)
	}
	return snap
}

func TestSystemCollector(t *testing.T) {
	col := collector.NewSystemCollector()
	snap := createTestSnapshot(t)

	if err := col.Collect(context.Background(), snap); err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if snap.Runtime() == nil {
		t.Fatalf("expected Runtime to be populated")
	}
	if snap.Runtime().OS == "" {
		t.Errorf("expected OS to be non-empty")
	}
	if snap.Runtime().Architecture == "" {
		t.Errorf("expected Architecture to be non-empty")
	}
	if snap.Runtime().State != domain.StateKnown {
		t.Errorf("expected StateKnown, got: %s", snap.Runtime().State)
	}
}

func TestRuntimeCollector(t *testing.T) {
	col := collector.NewRuntimeCollector()
	snap := createTestSnapshot(t)

	if err := col.Collect(context.Background(), snap); err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if snap.Runtime() == nil {
		t.Fatalf("expected Runtime to be populated")
	}
	// On this system, Go is installed
	if snap.Runtime().Runtimes == nil {
		t.Fatalf("expected Runtimes map to be non-nil")
	}
	if ver, ok := snap.Runtime().Runtimes["go"]; ok && ver == "" {
		t.Errorf("expected go version to be non-empty string if detected")
	}
}

func TestConfigCollector_Redaction(t *testing.T) {
	_ = os.Setenv("MY_APP_SECRET_TOKEN", "super-secret-value")
	_ = os.Setenv("MY_APP_NORMAL_VAR", "regular-value")
	defer os.Unsetenv("MY_APP_SECRET_TOKEN")
	defer os.Unsetenv("MY_APP_NORMAL_VAR")

	col := collector.NewConfigCollector()
	snap := createTestSnapshot(t)

	if err := col.Collect(context.Background(), snap); err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	cfg := snap.Configuration()
	if cfg == nil || cfg.Variables == nil {
		t.Fatalf("expected Configuration.Variables to be populated")
	}

	secretVar, ok := cfg.Variables["MY_APP_SECRET_TOKEN"]
	if !ok {
		t.Fatalf("expected MY_APP_SECRET_TOKEN in configuration")
	}
	if !secretVar.Sensitive || secretVar.Value != "[REDACTED]" || secretVar.State != domain.StateRedacted {
		t.Errorf("expected secret to be redacted, got: %+v", secretVar)
	}

	normalVar, ok := cfg.Variables["MY_APP_NORMAL_VAR"]
	if !ok {
		t.Fatalf("expected MY_APP_NORMAL_VAR in configuration")
	}
	if normalVar.Sensitive || normalVar.Value != "regular-value" || normalVar.State != domain.StateKnown {
		t.Errorf("expected normal var to be unredacted, got: %+v", normalVar)
	}
}

