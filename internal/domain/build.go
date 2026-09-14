package domain

type Build struct {
	ID string
	ToolChain map[string]string
	Dependencies map[string]string
}

func NewBuild(id string, toolChain, dependencies map[string]string) (*Build, error) {
	return &Build{
		ID: id,
		ToolChain: toolChain,
		Dependencies: dependencies,
	}, nil
}