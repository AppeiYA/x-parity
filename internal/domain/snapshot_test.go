package domain_test

import (
	"testing"
	"time"

	"github.com/AppeiYA/x-parity/internal/domain"
)

func TestNewSnapshot(t *testing.T) {
	app, err := domain.NewApplication("myapp", "git@github.com:example/myapp.git")
	if err != nil {
		t.Fatalf("unexpected error creating app: %v", err)
	}

	env, err := domain.NewEnvironment("production", "cloud")
	if err != nil {
		t.Fatalf("unexpected error creating env: %v", err)
	}

	now := time.Now().UTC()

	t.Run("successful creation", func(t *testing.T) {
		snap, err := domain.NewSnapshot("1.0.0", now, app, env)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if snap.Version() != "1.0.0" {
			t.Errorf("expected version 1.0.0, got: %s", snap.Version())
		}
		if snap.CapturedAt() != now {
			t.Errorf("expected capturedAt %v, got: %v", now, snap.CapturedAt())
		}
		if snap.Application() != app {
			t.Errorf("expected application %v, got: %v", app, snap.Application())
		}
		if snap.Environment() != env {
			t.Errorf("expected environment %v, got: %v", env, snap.Environment())
		}
		if snap.Dependencies() == nil {
			t.Errorf("expected dependencies to be non-nil empty slice")
		}
		if len(snap.Dependencies()) != 0 {
			t.Errorf("expected 0 dependencies, got: %d", len(snap.Dependencies()))
		}
	})

	t.Run("empty version fails", func(t *testing.T) {
		_, err := domain.NewSnapshot("", now, app, env)
		if err != domain.ErrEmptySnapshotVersion {
			t.Errorf("expected ErrEmptySnapshotVersion, got: %v", err)
		}
	})

	t.Run("nil application fails", func(t *testing.T) {
		_, err := domain.NewSnapshot("1.0.0", now, nil, env)
		if err != domain.ErrNilSnapshotApplication {
			t.Errorf("expected ErrNilSnapshotApplication, got: %v", err)
		}
	})

	t.Run("nil environment fails", func(t *testing.T) {
		_, err := domain.NewSnapshot("1.0.0", now, app, nil)
		if err != domain.ErrNilSnapshotEnvironment {
			t.Errorf("expected ErrNilSnapshotEnvironment, got: %v", err)
		}
	})
}

func TestSnapshotSetters(t *testing.T) {
	app, _ := domain.NewApplication("myapp", "")
	env, _ := domain.NewEnvironment("prod", "cloud")
	snap, err := domain.NewSnapshot("1.0.0", time.Now().UTC(), app, env)
	if err != nil {
		t.Fatalf("failed to create snapshot: %v", err)
	}

	// SetSource
	src, _ := domain.NewSource("github", "org/repo", "abc1234", "main", false, domain.StateKnown)
	snap.SetSource(src)
	if snap.Source() != src {
		t.Errorf("expected source %v, got: %v", src, snap.Source())
	}

	// SetBuild
	bld, _ := domain.NewBuild("build-1", map[string]string{"go": "1.25"}, map[string]string{})
	snap.SetBuild(bld)
	if snap.Build() != bld {
		t.Errorf("expected build %v, got: %v", bld, snap.Build())
	}

	// SetRuntime
	rt, _ := domain.NewRuntime("linux", "amd64", "6.5", "host-1", map[string]string{"go": "1.25"}, domain.StateKnown)
	snap.SetRuntime(rt)
	if snap.Runtime() != rt {
		t.Errorf("expected runtime %v, got: %v", rt, snap.Runtime())
	}

	// SetConfiguration
	cfgVal, _ := domain.NewConfigValue("8080", false, domain.StateKnown)
	cfg, _ := domain.NewConfiguration(map[string]*domain.ConfigValue{"PORT": cfgVal})
	snap.SetConfiguration(cfg)
	if snap.Configuration() != cfg {
		t.Errorf("expected configuration %v, got: %v", cfg, snap.Configuration())
	}

	// AddDependency
	dep1, _ := domain.NewDependency("postgres", "database", "15", "localhost", "5432")
	dep2, _ := domain.NewDependency("redis", "cache", "7", "localhost", "6379")
	snap.AddDependency(dep1)
	snap.AddDependency(dep2)
	if len(snap.Dependencies()) != 2 {
		t.Fatalf("expected 2 dependencies, got: %d", len(snap.Dependencies()))
	}
	if snap.Dependencies()[0] != dep1 || snap.Dependencies()[1] != dep2 {
		t.Errorf("dependencies do not match expected items")
	}

	// SetDependencies
	snap.SetDependencies([]*domain.Dependency{dep1})
	if len(snap.Dependencies()) != 1 || snap.Dependencies()[0] != dep1 {
		t.Errorf("SetDependencies did not overwrite dependencies slice")
	}

	// SetDeployment
	dep, _ := domain.NewDeployment("k8s", "myapp:latest", "sha256:123")
	snap.SetDeployment(dep)
	if snap.Deployment() != dep {
		t.Errorf("expected deployment %v, got: %v", dep, snap.Deployment())
	}

	// SetApplication & SetEnvironment
	newApp, _ := domain.NewApplication("myapp2", "git@github.com:example/myapp2.git")
	newEnv, _ := domain.NewEnvironment("staging", "cloud")
	snap.SetApplication(newApp)
	snap.SetEnvironment(newEnv)
	if snap.Application() != newApp {
		t.Errorf("expected application %v, got: %v", newApp, snap.Application())
	}
	if snap.Environment() != newEnv {
		t.Errorf("expected environment %v, got: %v", newEnv, snap.Environment())
	}
}

