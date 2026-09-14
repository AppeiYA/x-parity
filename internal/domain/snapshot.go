package domain

import (
	"fmt"
	"time"
)

type Snapshot struct {
	version string
	capturedAt time.Time
	application *Application
	environment *Environment
	source *Source
	build *Build
	runtime *Runtime
	configuration *Configuration
	dependencies []*Dependency
	deployment *Deployment
}

var (
	ErrEmptySnapshotVersion   = fmt.Errorf("snapshot version cannot be empty")
	ErrNilSnapshotApplication = fmt.Errorf("snapshot application cannot be nil")
	ErrNilSnapshotEnvironment = fmt.Errorf("snapshot environment cannot be nil")
)

func NewSnapshot(version string, capturedAt time.Time, application *Application, environment *Environment) (*Snapshot, error) {
	if version == "" {
		return nil, ErrEmptySnapshotVersion
	}
	if application == nil {
		return nil, ErrNilSnapshotApplication
	}
	if environment == nil {
		return nil, ErrNilSnapshotEnvironment
	}
	return &Snapshot{
		version:      version,
		capturedAt:   capturedAt,
		application:  application,
		environment:  environment,
		dependencies: make([]*Dependency, 0),
	}, nil
}

func (s *Snapshot) Version() string {
	return s.version
}

func (s *Snapshot) CapturedAt() time.Time {
	return s.capturedAt
}

func (s *Snapshot) Application() *Application {
	return s.application
}

func (s *Snapshot) Environment() *Environment {
	return s.environment
}

func (s *Snapshot) Source() *Source {
	return s.source
}

func (s *Snapshot) Build() *Build {
	return s.build
}

func (s *Snapshot) Runtime() *Runtime {
	return s.runtime
}

func (s *Snapshot) Configuration() *Configuration {
	return s.configuration
}

func (s *Snapshot) Dependencies() []*Dependency {
	return s.dependencies
}

func (s *Snapshot) Deployment() *Deployment {
	return s.deployment
}

func (s *Snapshot) SetApplication(application *Application) {
	s.application = application
}

func (s *Snapshot) SetEnvironment(environment *Environment) {
	s.environment = environment
}

func (s *Snapshot) SetSource(source *Source) {
	s.source = source
}

func (s *Snapshot) SetBuild(build *Build) {
	s.build = build
}

func (s *Snapshot) SetRuntime(runtime *Runtime) {
	s.runtime = runtime
}

func (s *Snapshot) SetConfiguration(configuration *Configuration) {
	s.configuration = configuration
}

func (s *Snapshot) SetDependencies(dependencies []*Dependency) {
	s.dependencies = dependencies
}

func (s *Snapshot) AddDependency(dependency *Dependency) {
	s.dependencies = append(s.dependencies, dependency)
}

func (s *Snapshot) SetDeployment(deployment *Deployment) {
	s.deployment = deployment
}