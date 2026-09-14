package domain

import "fmt"

type Dependency struct {
	Name string 
	Type string
	Version string
	Host string
	Port string
}

var (
	ErrEmptyDependencyName = fmt.Errorf("dependency name cannot be empty")
	ErrEmptyDependencyType = fmt.Errorf("dependency type cannot be empty")
)

func NewDependency(name, depType, version, host, port string) (*Dependency, error) {
	if name == "" {
		return nil, ErrEmptyDependencyName
	}

	if depType == "" {
		return nil, ErrEmptyDependencyType
	}

	return &Dependency{
		Name: name,
		Type: depType,
		Version: version,
		Host: host,
		Port: port,
	}, nil
}