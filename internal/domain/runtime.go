package domain

import "fmt"

type Runtime struct {
	OS string
	Architecture string
	Kernel string
	Hostname string
	Virtualization string
	Runtimes map[string]string
	State EvidenceState
}

var (
	ErrRuntimeOSCannotBeEmpty = fmt.Errorf("runtime OS cannot be empty")
	ErrRuntimeArchitectureCannotBeEmpty = fmt.Errorf("runtime architecture cannot be empty")
)

func NewRuntime(os, architecture, kernel, hostname string, runtimes map[string]string, state EvidenceState) (*Runtime, error) {
	if os == "" {
		return nil, ErrRuntimeOSCannotBeEmpty
	}
	if architecture == "" {
		return nil, ErrRuntimeArchitectureCannotBeEmpty
	}
	if !state.IsValid() {
		return nil, ErrInvalidEvidenceState
	}
	return &Runtime{
		OS: os,
		Architecture: architecture,
		Kernel: kernel,
		Hostname: hostname,
		Runtimes: runtimes,
		State: state,
	}, nil
}