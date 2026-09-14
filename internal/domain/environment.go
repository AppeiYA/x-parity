package domain

import "fmt"

type Environment struct {
	Name string
	Type string
}

var (
	ErrEmptyEnvironmentName = fmt.Errorf("environment name cannot be empty")
	ErrEmptyEnvironmentType = fmt.Errorf("environment type cannot be empty")
)

func NewEnvironment(name, envType string) (*Environment, error) {
	if name == "" {
		return nil, ErrEmptyEnvironmentName
	}

	if envType == "" {
		return nil, ErrEmptyEnvironmentType
	}
	return &Environment{
		Name: name,
		Type: envType,
	}, nil
}