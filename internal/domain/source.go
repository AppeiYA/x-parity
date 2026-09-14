package domain

import "fmt"

type Source struct {
	Provider string
	Repository string
	Commit string
	Branch string
	Dirty bool
	State EvidenceState
}

var (
	ErrInvalidEvidenceState = fmt.Errorf("invalid evidence state")
)

func NewSource(provider, repository, commit, branch string, dirty bool, state EvidenceState) (*Source, error) {
	if !state.IsValid() {
		return nil, ErrInvalidEvidenceState
	}
	return &Source{
		Provider: provider,
		Repository: repository,
		Commit: commit,
		Branch: branch,
		Dirty: dirty,
		State: state,
	}, nil
}