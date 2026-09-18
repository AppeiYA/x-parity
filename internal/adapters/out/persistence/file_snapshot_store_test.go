package persistence_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/AppeiYA/x-parity/internal/adapters/out/persistence"
	"github.com/AppeiYA/x-parity/internal/domain"
)

func TestFileSnapshotStore_SaveAndLoad(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snapshot-store-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "snapshots", "test_snap.json")

	app, _ := domain.NewApplication("myapp", "git@github.com:org/myapp.git")
	env, _ := domain.NewEnvironment("production", "cloud")
	now := time.Now().UTC().Truncate(time.Millisecond)

	snap, err := domain.NewSnapshot("1.0.0", now, app, env)
	if err != nil {
		t.Fatalf("failed to create snapshot: %v", err)
	}

	src, _ := domain.NewSource("github", "org/myapp", "commit123", "main", false, domain.StateKnown)
	snap.SetSource(src)

	rt, _ := domain.NewRuntime("linux", "amd64", "6.5.0", "host-prod", map[string]string{"go": "1.25", "node": "20.0"}, domain.StateKnown)
	snap.SetRuntime(rt)

	val, _ := domain.NewConfigValue("prod-secret", true, domain.StateRedacted)
	cfg, _ := domain.NewConfiguration(map[string]*domain.ConfigValue{"API_KEY": val})
	snap.SetConfiguration(cfg)

	dep, _ := domain.NewDependency("postgres", "db", "16", "db.internal", "5432")
	snap.AddDependency(dep)

	store := persistence.NewFileSnapshotStore()

	ctx := context.Background()
	if err := store.Save(ctx, filePath, snap); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file was written and non-empty
	info, err := os.Stat(filePath)
	if err != nil || info.Size() == 0 {
		t.Fatalf("file does not exist or is empty")
	}

	loaded, err := store.Load(ctx, filePath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Version() != snap.Version() {
		t.Errorf("expected version %s, got %s", snap.Version(), loaded.Version())
	}
	if loaded.Application().Name != "myapp" || loaded.Application().Repository != "git@github.com:org/myapp.git" {
		t.Errorf("application mismatch: %v", loaded.Application())
	}
	if loaded.Environment().Name != "production" {
		t.Errorf("environment mismatch: %v", loaded.Environment())
	}
	if loaded.Source() == nil || loaded.Source().Commit != "commit123" {
		t.Errorf("source mismatch: %v", loaded.Source())
	}
	if loaded.Runtime() == nil || loaded.Runtime().OS != "linux" || loaded.Runtime().Runtimes["go"] != "1.25" {
		t.Errorf("runtime mismatch: %v", loaded.Runtime())
	}
	if loaded.Configuration() == nil || loaded.Configuration().Variables["API_KEY"].Sensitive != true {
		t.Errorf("configuration mismatch: %v", loaded.Configuration())
	}
	if len(loaded.Dependencies()) != 1 || loaded.Dependencies()[0].Name != "postgres" {
		t.Errorf("dependencies mismatch: %v", loaded.Dependencies())
	}
}

