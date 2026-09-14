package domain

import "fmt"

type Application struct {
	Name string
	Repository string
}

var (
	ErrEmptyApplicationName = fmt.Errorf("application name cannot be empty")
)

func NewApplication(name, repository string) (*Application, error) {
	if name == "" {
		return nil, ErrEmptyApplicationName
	}
	return &Application{
		Name: name,
		Repository: repository,
	}, nil
}