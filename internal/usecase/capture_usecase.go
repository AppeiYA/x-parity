package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/AppeiYA/x-parity/internal/domain"
	portin "github.com/AppeiYA/x-parity/internal/ports/in"
	portout "github.com/AppeiYA/x-parity/internal/ports/out"
)

type CaptureUsecase struct {
	collectors []portout.CollectorInt
	store      portout.SnapshotStore
}

func NewCaptureUsecase(collectors []portout.CollectorInt, store portout.SnapshotStore) *CaptureUsecase {
	return &CaptureUsecase{
		collectors: collectors,
		store:      store,
	}
}

func (c *CaptureUsecase) Execute(ctx context.Context, cmd portin.CaptureCommand) (*domain.Snapshot, error) {
	app, err := domain.NewApplication(cmd.AppName, "")
	if err != nil {
		return nil, err
	}

	env, err := domain.NewEnvironment(
		cmd.Environment,
		"captured",
	)
	if err != nil {
		return nil, err
	}

	snapshot, err := domain.NewSnapshot(
		"1.0.0",
		time.Now().UTC(),
		app,
		env,
	)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, col := range c.collectors {
		wg.Add(1)
		go func(collector portout.CollectorInt) {
			defer wg.Done()

			if err := collector.Collect(ctx, snapshot); err != nil {
				mu.Lock()
				switch collector.Name() {
				case "source":
					src, _ := domain.NewSource("", "", "", "", false, domain.StateUnavailable)
					snapshot.SetSource(src)
				case "runtime":
					rt, _ := domain.NewRuntime("unknown", "unknown", "", "", nil, domain.StateUnavailable)
					snapshot.SetRuntime(rt)
				case "config", "configuration":
					cfg, _ := domain.NewConfiguration(make(map[string]*domain.ConfigValue))
					snapshot.SetConfiguration(cfg)
				case "build":
					bld, _ := domain.NewBuild("", nil, nil)
					snapshot.SetBuild(bld)
				case "deployment":
					dep, _ := domain.NewDeployment("", "", "")
					snapshot.SetDeployment(dep)
				}
				mu.Unlock()
			}
		}(col)
	}

	wg.Wait()

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if cmd.OutputPath != "" && c.store != nil {
		if err := c.store.Save(ctx, cmd.OutputPath, snapshot); err != nil {
			return nil, err
		}
	}

	return snapshot, nil
}
