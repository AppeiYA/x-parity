package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/AppeiYA/x-parity/internal/domain"
	portout "github.com/AppeiYA/x-parity/internal/ports/out"
)

var _ portout.SnapshotStore = (*FileSnapshotStore)(nil)

type FileSnapshotStore struct{}

func NewFileSnapshotStore() *FileSnapshotStore {
	return &FileSnapshotStore{}
}

type snapshotDTO struct {
	Version       string            `json:"version"`
	CapturedAt    time.Time         `json:"captured_at"`
	Application   *applicationDTO   `json:"application,omitempty"`
	Environment   *environmentDTO   `json:"environment,omitempty"`
	Source        *sourceDTO        `json:"source,omitempty"`
	Build         *buildDTO         `json:"build,omitempty"`
	Runtime       *runtimeDTO       `json:"runtime,omitempty"`
	Configuration *configurationDTO `json:"configuration,omitempty"`
	Dependencies  []*dependencyDTO  `json:"dependencies,omitempty"`
	Deployment    *deploymentDTO    `json:"deployment,omitempty"`
}

type applicationDTO struct {
	Name       string `json:"name"`
	Repository string `json:"repository,omitempty"`
}

type environmentDTO struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type sourceDTO struct {
	Provider   string `json:"provider,omitempty"`
	Repository string `json:"repository,omitempty"`
	Commit     string `json:"commit,omitempty"`
	Branch     string `json:"branch,omitempty"`
	Dirty      bool   `json:"dirty"`
	State      string `json:"state"`
}

type buildDTO struct {
	ID           string            `json:"id,omitempty"`
	ToolChain    map[string]string `json:"tool_chain,omitempty"`
	Dependencies map[string]string `json:"dependencies,omitempty"`
}

type runtimeDTO struct {
	OS             string            `json:"os"`
	Architecture   string            `json:"architecture"`
	Kernel         string            `json:"kernel,omitempty"`
	Hostname       string            `json:"hostname,omitempty"`
	Virtualization string            `json:"virtualization,omitempty"`
	Runtimes       map[string]string `json:"runtimes,omitempty"`
	State          string            `json:"state"`
}

type configurationDTO struct {
	Variables map[string]*configValueDTO `json:"variables,omitempty"`
}

type configValueDTO struct {
	State     string `json:"state"`
	Value     string `json:"value"`
	Sensitive bool   `json:"sensitive"`
}

type dependencyDTO struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Version string `json:"version,omitempty"`
	Host    string `json:"host,omitempty"`
	Port    string `json:"port,omitempty"`
}

type deploymentDTO struct {
	Platform    string `json:"platform,omitempty"`
	Image       string `json:"image,omitempty"`
	ImageDigest string `json:"image_digest,omitempty"`
}

func (s *FileSnapshotStore) Save(ctx context.Context, path string, snapshot *domain.Snapshot) error {
	if snapshot == nil {
		return fmt.Errorf("cannot save nil snapshot")
	}

	dto := toDTO(snapshot)
	data, err := json.MarshalIndent(dto, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	tmpFile := path + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write temp snapshot: %w", err)
	}

	if err := os.Rename(tmpFile, path); err != nil {
		return fmt.Errorf("failed to rename temp snapshot to %s: %w", path, err)
	}

	return nil
}

func (s *FileSnapshotStore) Load(ctx context.Context, path string) (*domain.Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read snapshot file: %w", err)
	}

	var dto snapshotDTO
	if err := json.Unmarshal(data, &dto); err != nil {
		return nil, fmt.Errorf("failed to unmarshal snapshot JSON: %w", err)
	}

	return fromDTO(&dto)
}

func toDTO(snap *domain.Snapshot) *snapshotDTO {
	dto := &snapshotDTO{
		Version:    snap.Version(),
		CapturedAt: snap.CapturedAt(),
	}

	if app := snap.Application(); app != nil {
		dto.Application = &applicationDTO{
			Name:       app.Name,
			Repository: app.Repository,
		}
	}

	if env := snap.Environment(); env != nil {
		dto.Environment = &environmentDTO{
			Name: env.Name,
			Type: env.Type,
		}
	}

	if src := snap.Source(); src != nil {
		dto.Source = &sourceDTO{
			Provider:   src.Provider,
			Repository: src.Repository,
			Commit:     src.Commit,
			Branch:     src.Branch,
			Dirty:      src.Dirty,
			State:      string(src.State),
		}
	}

	if bld := snap.Build(); bld != nil {
		dto.Build = &buildDTO{
			ID:           bld.ID,
			ToolChain:    bld.ToolChain,
			Dependencies: bld.Dependencies,
		}
	}

	if rt := snap.Runtime(); rt != nil {
		dto.Runtime = &runtimeDTO{
			OS:             rt.OS,
			Architecture:   rt.Architecture,
			Kernel:         rt.Kernel,
			Hostname:       rt.Hostname,
			Virtualization: rt.Virtualization,
			Runtimes:       rt.Runtimes,
			State:          string(rt.State),
		}
	}

	if cfg := snap.Configuration(); cfg != nil {
		dto.Configuration = &configurationDTO{
			Variables: make(map[string]*configValueDTO),
		}
		for k, v := range cfg.Variables {
			if v != nil {
				dto.Configuration.Variables[k] = &configValueDTO{
					State:     string(v.State),
					Value:     v.Value,
					Sensitive: v.Sensitive,
				}
			}
		}
	}

	if deps := snap.Dependencies(); len(deps) > 0 {
		dto.Dependencies = make([]*dependencyDTO, 0, len(deps))
		for _, d := range deps {
			if d != nil {
				dto.Dependencies = append(dto.Dependencies, &dependencyDTO{
					Name:    d.Name,
					Type:    d.Type,
					Version: d.Version,
					Host:    d.Host,
					Port:    d.Port,
				})
			}
		}
	}

	if dep := snap.Deployment(); dep != nil {
		dto.Deployment = &deploymentDTO{
			Platform:    dep.Platform,
			Image:       dep.Image,
			ImageDigest: dep.ImageDigest,
		}
	}

	return dto
}

func fromDTO(dto *snapshotDTO) (*domain.Snapshot, error) {
	var app *domain.Application
	if dto.Application != nil {
		var err error
		app, err = domain.NewApplication(dto.Application.Name, dto.Application.Repository)
		if err != nil {
			return nil, fmt.Errorf("invalid application in snapshot: %w", err)
		}
	} else {
		return nil, fmt.Errorf("snapshot is missing application")
	}

	var env *domain.Environment
	if dto.Environment != nil {
		var err error
		env, err = domain.NewEnvironment(dto.Environment.Name, dto.Environment.Type)
		if err != nil {
			return nil, fmt.Errorf("invalid environment in snapshot: %w", err)
		}
	} else {
		return nil, fmt.Errorf("snapshot is missing environment")
	}

	snap, err := domain.NewSnapshot(dto.Version, dto.CapturedAt, app, env)
	if err != nil {
		return nil, fmt.Errorf("failed to construct snapshot: %w", err)
	}

	if dto.Source != nil {
		st := domain.EvidenceState(dto.Source.State)
		if !st.IsValid() {
			st = domain.StateKnown
		}
		src, err := domain.NewSource(dto.Source.Provider, dto.Source.Repository, dto.Source.Commit, dto.Source.Branch, dto.Source.Dirty, st)
		if err == nil {
			snap.SetSource(src)
		}
	}

	if dto.Build != nil {
		bld, err := domain.NewBuild(dto.Build.ID, dto.Build.ToolChain, dto.Build.Dependencies)
		if err == nil {
			snap.SetBuild(bld)
		}
	}

	if dto.Runtime != nil {
		st := domain.EvidenceState(dto.Runtime.State)
		if !st.IsValid() {
			st = domain.StateKnown
		}
		rt, err := domain.NewRuntime(dto.Runtime.OS, dto.Runtime.Architecture, dto.Runtime.Kernel, dto.Runtime.Hostname, dto.Runtime.Runtimes, st)
		if err == nil {
			rt.Virtualization = dto.Runtime.Virtualization
			snap.SetRuntime(rt)
		}
	}

	if dto.Configuration != nil {
		vars := make(map[string]*domain.ConfigValue)
		for k, v := range dto.Configuration.Variables {
			if v != nil {
				st := domain.EvidenceState(v.State)
				if !st.IsValid() {
					st = domain.StateKnown
				}
				cv, err := domain.NewConfigValue(v.Value, v.Sensitive, st)
				if err == nil {
					vars[k] = cv
				}
			}
		}
		cfg, err := domain.NewConfiguration(vars)
		if err == nil {
			snap.SetConfiguration(cfg)
		}
	}

	if len(dto.Dependencies) > 0 {
		for _, d := range dto.Dependencies {
			if d != nil {
				dep, err := domain.NewDependency(d.Name, d.Type, d.Version, d.Host, d.Port)
				if err == nil {
					snap.AddDependency(dep)
				}
			}
		}
	}

	if dto.Deployment != nil {
		dep, err := domain.NewDeployment(dto.Deployment.Platform, dto.Deployment.Image, dto.Deployment.ImageDigest)
		if err == nil {
			snap.SetDeployment(dep)
		}
	}

	return snap, nil
}