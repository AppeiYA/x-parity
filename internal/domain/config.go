package domain

import "fmt"

type ConfigValue struct {
	State EvidenceState
	Value string
	Sensitive bool
}

type Configuration struct {
	Variables map[string]*ConfigValue
}

var (
	ErrEmptyConfigurationValues = fmt.Errorf("configuration values cannot be empty")
)

func NewConfigValue(value string, sensitive bool, state EvidenceState) (*ConfigValue, error) {
	if !state.IsValid() {
		return nil, ErrInvalidEvidenceState
	}
	return &ConfigValue{
		State: state,
		Value: value,
		Sensitive: sensitive,
	}, nil
}

func NewConfiguration(variables map[string]*ConfigValue) (*Configuration, error) {
	if variables == nil {
		return nil, ErrEmptyConfigurationValues
	}
	return &Configuration{
		Variables: variables,
	}, nil
}