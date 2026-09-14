package usecase_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/AppeiYA/x-parity/internal/domain"
	"github.com/AppeiYA/x-parity/internal/usecase"
)

func TestCompareUsecase_Execute(t *testing.T) {
	app, _ := domain.NewApplication("myapp", "")
	envLocal, _ := domain.NewEnvironment("local", "dev")
	envRemote, _ := domain.NewEnvironment("prod", "cloud")

	snapLocal, _ := domain.NewSnapshot("1.0.0", time.Now().UTC(), app, envLocal)
	snapRemote, _ := domain.NewSnapshot("1.0.0", time.Now().UTC(), app, envRemote)

	rtLocal, _ := domain.NewRuntime("linux", "amd64", "6.5", "host-1", map[string]string{"go": "1.25"}, domain.StateKnown)
	rtRemote, _ := domain.NewRuntime("darwin", "arm64", "6.5", "host-2", map[string]string{"go": "1.25"}, domain.StateKnown)
	snapLocal.SetRuntime(rtLocal)
	snapRemote.SetRuntime(rtRemote)

	store := newMockStore()
	_ = store.Save(context.Background(), "/path/local.json", snapLocal)
	_ = store.Save(context.Background(), "/path/remote.json", snapRemote)

	diffEngine := domain.NewDiffEngine()
	uc := usecase.NewCompareUsecase(store, diffEngine)

	t.Run("successful comparison", func(t *testing.T) {
		diffs, err := uc.Execute(context.Background(), "/path/local.json", "/path/remote.json")
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if len(diffs) == 0 {
			t.Fatalf("expected differences due to OS divergence, got 0")
		}
		if diffs[0].Path() != "runtime.os" {
			t.Errorf("expected runtime.os diff, got %s", diffs[0].Path())
		}
	})

	t.Run("empty local path returns error", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), "", "/path/remote.json")
		if !errors.Is(err, usecase.ErrEmptyLocalPath) {
			t.Errorf("expected ErrEmptyLocalPath, got: %v", err)
		}
	})

	t.Run("empty remote path returns error", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), "/path/local.json", "")
		if !errors.Is(err, usecase.ErrEmptyRemotePath) {
			t.Errorf("expected ErrEmptyRemotePath, got: %v", err)
		}
	})

	t.Run("local path not found returns wrapped error", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), "/path/nonexistent.json", "/path/remote.json")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to load local snapshot") {
			t.Errorf("expected error message to contain context, got: %v", err)
		}
	})

	t.Run("remote path not found returns wrapped error", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), "/path/local.json", "/path/nonexistent.json")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to load remote snapshot") {
			t.Errorf("expected error message to contain context, got: %v", err)
		}
	})
}

