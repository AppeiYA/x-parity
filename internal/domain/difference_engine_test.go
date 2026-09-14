package domain_test

import (
	"testing"
	"time"

	"github.com/AppeiYA/x-parity/internal/domain"
)

func createBaseSnapshot(t *testing.T, name, envName string) *domain.Snapshot {
	t.Helper()
	app, err := domain.NewApplication(name, "")
	if err != nil {
		t.Fatalf("failed to create application: %v", err)
	}
	env, err := domain.NewEnvironment(envName, "cloud")
	if err != nil {
		t.Fatalf("failed to create environment: %v", err)
	}
	snap, err := domain.NewSnapshot("1.0.0", time.Now().UTC(), app, env)
	if err != nil {
		t.Fatalf("failed to create snapshot: %v", err)
	}
	return snap
}

func TestDiffEngine_Compare_NilSafety(t *testing.T) {
	engine := domain.NewDiffEngine()

	t.Run("partial snapshots with nil subsystems do not panic", func(t *testing.T) {
		local := createBaseSnapshot(t, "app", "local")
		remote := createBaseSnapshot(t, "app", "remote")

		diffs, err := engine.Compare(*local, *remote)
		if err != nil {
			t.Fatalf("expected no error for empty subsystems, got: %v", err)
		}
		if len(diffs) != 0 {
			t.Errorf("expected 0 diffs for identical empty snapshots, got: %d", len(diffs))
		}
	})

	t.Run("one snapshot missing runtime does not panic", func(t *testing.T) {
		local := createBaseSnapshot(t, "app", "local")
		rt, err := domain.NewRuntime("linux", "amd64", "6.5", "host-1", map[string]string{"go": "1.25"}, domain.StateKnown)
		if err != nil {
			t.Fatalf("failed to create runtime: %v", err)
		}
		local.SetRuntime(rt)

		remote := createBaseSnapshot(t, "app", "remote")

		diffs, err := engine.Compare(*local, *remote)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Since remote has no runtime, it should safely skip without panic
		if len(diffs) != 0 {
			t.Errorf("expected 0 diffs, got: %d", len(diffs))
		}
	})

	t.Run("one snapshot missing configuration does not panic", func(t *testing.T) {
		local := createBaseSnapshot(t, "app", "local")
		remote := createBaseSnapshot(t, "app", "remote")

		cfgVal, _ := domain.NewConfigValue("8080", false, domain.StateKnown)
		cfg, _ := domain.NewConfiguration(map[string]*domain.ConfigValue{"PORT": cfgVal})
		remote.SetConfiguration(cfg)

		diffs, err := engine.Compare(*local, *remote)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(diffs) != 0 {
			t.Errorf("expected 0 diffs, got: %d", len(diffs))
		}
	})
}

func TestDiffEngine_Compare_TelemetryDifferences(t *testing.T) {
	engine := domain.NewDiffEngine()

	local := createBaseSnapshot(t, "app", "local")
	remote := createBaseSnapshot(t, "app", "prod")

	// Runtimes
	rtLocal, _ := domain.NewRuntime("linux", "amd64", "6.5", "host-1", map[string]string{"go": "1.25.1", "node": "20.1.0"}, domain.StateKnown)
	rtRemote, _ := domain.NewRuntime("darwin", "arm64", "23.0", "host-2", map[string]string{"go": "1.24.0", "python": "3.11"}, domain.StateKnown)
	local.SetRuntime(rtLocal)
	remote.SetRuntime(rtRemote)

	// Configurations
	val1, _ := domain.NewConfigValue("debug", false, domain.StateKnown)
	cfgLocal, _ := domain.NewConfiguration(map[string]*domain.ConfigValue{"LOG_LEVEL": val1})
	val2, _ := domain.NewConfigValue("info", false, domain.StateKnown)
	cfgRemote, _ := domain.NewConfiguration(map[string]*domain.ConfigValue{"LOG_LEVEL": val2})
	local.SetConfiguration(cfgLocal)
	remote.SetConfiguration(cfgRemote)

	// Sources
	srcLocal, _ := domain.NewSource("github", "org/repo", "sha-local", "feature", false, domain.StateKnown)
	srcRemote, _ := domain.NewSource("github", "org/repo", "sha-remote", "main", false, domain.StateKnown)
	local.SetSource(srcLocal)
	remote.SetSource(srcRemote)

	diffs, err := engine.Compare(*local, *remote)
	if err != nil {
		t.Fatalf("unexpected compare error: %v", err)
	}

	if len(diffs) == 0 {
		t.Fatalf("expected multiple differences, got 0")
	}

	hasOSDiff := false
	hasConfigDiff := false
	hasSourceDiff := false

	for _, d := range diffs {
		switch d.Path() {
		case "runtime.os":
			hasOSDiff = true
		case "configuration.variables.LOG_LEVEL":
			hasConfigDiff = true
		case "source.commit":
			hasSourceDiff = true
		}
	}

	if !hasOSDiff {
		t.Errorf("expected runtime.os difference")
	}
	if !hasConfigDiff {
		t.Errorf("expected configuration.variables.LOG_LEVEL difference")
	}
	if !hasSourceDiff {
		t.Errorf("expected source.commit difference")
	}
}
