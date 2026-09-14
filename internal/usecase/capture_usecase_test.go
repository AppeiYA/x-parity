package usecase_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/AppeiYA/x-parity/internal/domain"
	portin "github.com/AppeiYA/x-parity/internal/ports/in"
	portout "github.com/AppeiYA/x-parity/internal/ports/out"
	"github.com/AppeiYA/x-parity/internal/usecase"
)

type mockCollector struct {
	name    string
	collect func(ctx context.Context, snapshot *domain.Snapshot) error
}

func (m *mockCollector) Name() string {
	return m.name
}

func (m *mockCollector) Collect(ctx context.Context, snapshot *domain.Snapshot) error {
	if m.collect != nil {
		return m.collect(ctx, snapshot)
	}
	return nil
}

type mockStore struct {
	mu    sync.Mutex
	saved map[string]*domain.Snapshot
}

func newMockStore() *mockStore {
	return &mockStore{
		saved: make(map[string]*domain.Snapshot),
	}
}

func (m *mockStore) Save(ctx context.Context, path string, snapshot *domain.Snapshot) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.saved[path] = snapshot
	return nil
}

func (m *mockStore) Load(ctx context.Context, path string) (*domain.Snapshot, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	snap, ok := m.saved[path]
	if !ok {
		return nil, errors.New("not found")
	}
	return snap, nil
}

func TestCaptureUsecase_Execute_Success(t *testing.T) {
	sourceCol := &mockCollector{
		name: "source",
		collect: func(ctx context.Context, snapshot *domain.Snapshot) error {
			src, err := domain.NewSource("github", "org/repo", "sha123", "main", false, domain.StateKnown)
			if err != nil {
				return err
			}
			snapshot.SetSource(src)
			return nil
		},
	}

	runtimeCol := &mockCollector{
		name: "runtime",
		collect: func(ctx context.Context, snapshot *domain.Snapshot) error {
			rt, err := domain.NewRuntime("linux", "arm64", "6.1", "host-a", map[string]string{"go": "1.25"}, domain.StateKnown)
			if err != nil {
				return err
			}
			snapshot.SetRuntime(rt)
			return nil
		},
	}

	store := newMockStore()
	uc := usecase.NewCaptureUsecase([]portout.CollectorInt{sourceCol, runtimeCol}, store)

	ctx := context.Background()
	cmd := portin.CaptureCommand{
		AppName:     "test-app",
		Environment: "staging",
		OutputPath:  "/tmp/snapshot.json",
	}

	snap, err := uc.Execute(ctx, cmd)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if snap.Application().Name != "test-app" {
		t.Errorf("expected app name test-app, got %s", snap.Application().Name)
	}
	if snap.Environment().Name != "staging" {
		t.Errorf("expected env name staging, got %s", snap.Environment().Name)
	}
	if snap.Source() == nil || snap.Source().Commit != "sha123" {
		t.Errorf("expected source with commit sha123, got %v", snap.Source())
	}
	if snap.Runtime() == nil || snap.Runtime().Architecture != "arm64" {
		t.Errorf("expected runtime arch arm64, got %v", snap.Runtime())
	}

	if savedSnap, ok := store.saved["/tmp/snapshot.json"]; !ok || savedSnap != snap {
		t.Errorf("expected snapshot to be saved to store at /tmp/snapshot.json")
	}
}

func TestCaptureUsecase_Execute_GracefulDegradation(t *testing.T) {
	sourceCol := &mockCollector{
		name: "source",
		collect: func(ctx context.Context, snapshot *domain.Snapshot) error {
			return errors.New("git command failed")
		},
	}

	runtimeCol := &mockCollector{
		name: "runtime",
		collect: func(ctx context.Context, snapshot *domain.Snapshot) error {
			return errors.New("uname command failed")
		},
	}

	uc := usecase.NewCaptureUsecase([]portout.CollectorInt{sourceCol, runtimeCol}, nil)

	ctx := context.Background()
	cmd := portin.CaptureCommand{
		AppName:     "test-app",
		Environment: "production",
	}

	snap, err := uc.Execute(ctx, cmd)
	if err != nil {
		t.Fatalf("expected graceful degradation without error, got: %v", err)
	}

	if snap.Source() == nil {
		t.Fatalf("expected non-nil fallback source")
	}
	if snap.Source().State != domain.StateUnavailable {
		t.Errorf("expected source state StateUnavailable, got: %v", snap.Source().State)
	}

	if snap.Runtime() == nil {
		t.Fatalf("expected non-nil fallback runtime")
	}
	if snap.Runtime().State != domain.StateUnavailable {
		t.Errorf("expected runtime state StateUnavailable, got: %v", snap.Runtime().State)
	}
}

func TestCaptureUsecase_Execute_EmptyAppName(t *testing.T) {
	uc := usecase.NewCaptureUsecase(nil, nil)
	cmd := portin.CaptureCommand{
		AppName:     "",
		Environment: "production",
	}

	_, err := uc.Execute(context.Background(), cmd)
	if err != domain.ErrEmptyApplicationName {
		t.Errorf("expected ErrEmptyApplicationName, got: %v", err)
	}
}

