package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/AppeiYA/x-parity/internal/domain"
	"github.com/AppeiYA/x-parity/internal/usecase"
)

func TestInspectUsecase(t *testing.T) {
	app, _ := domain.NewApplication("myapp", "")
	env, _ := domain.NewEnvironment("prod", "cloud")
	snap, _ := domain.NewSnapshot("1.0.0", time.Now().UTC(), app, env)

	store := newMockStore()
	_ = store.Save(context.Background(), "/path/snap.json", snap)

	uc := usecase.NewInspectUsecase(store)

	t.Run("successful inspection", func(t *testing.T) {
		loaded, err := uc.Execute(context.Background(), "/path/snap.json")
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if loaded.Application().Name != "myapp" {
			t.Errorf("expected myapp, got %s", loaded.Application().Name)
		}
	})

	t.Run("empty path fails", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), "")
		if err == nil {
			t.Errorf("expected error on empty path, got nil")
		}
	})
}

func TestDiagnoseUsecase(t *testing.T) {
	app, _ := domain.NewApplication("myapp", "")
	env, _ := domain.NewEnvironment("prod", "cloud")

	t.Run("healthy snapshot", func(t *testing.T) {
		snap, _ := domain.NewSnapshot("1.0.0", time.Now().UTC(), app, env)
		rt, _ := domain.NewRuntime("linux", "amd64", "6.5", "host", map[string]string{"go": "1.25"}, domain.StateKnown)
		snap.SetRuntime(rt)

		store := newMockStore()
		_ = store.Save(context.Background(), "/path/healthy.json", snap)

		uc := usecase.NewDiagnoseUsecase(store, domain.NewDiffEngine())
		report, err := uc.Execute(context.Background(), "/path/healthy.json")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(report.Hypotheses) != 0 {
			t.Errorf("expected 0 hypotheses for healthy snapshot, got %d", len(report.Hypotheses))
		}
	})

	t.Run("degraded snapshot with dirty git and missing config", func(t *testing.T) {
		snap, _ := domain.NewSnapshot("1.0.0", time.Now().UTC(), app, env)
		src, _ := domain.NewSource("github", "org/repo", "sha123", "main", true, domain.StateKnown)
		snap.SetSource(src)

		val, _ := domain.NewConfigValue("", false, domain.StateMissing)
		cfg, _ := domain.NewConfiguration(map[string]*domain.ConfigValue{"DB_HOST": val})
		snap.SetConfiguration(cfg)

		store := newMockStore()
		_ = store.Save(context.Background(), "/path/degraded.json", snap)

		uc := usecase.NewDiagnoseUsecase(store, domain.NewDiffEngine())
		report, err := uc.Execute(context.Background(), "/path/degraded.json")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(report.Hypotheses) < 2 {
			t.Errorf("expected at least 2 hypotheses, got %d", len(report.Hypotheses))
		}
	})
}

