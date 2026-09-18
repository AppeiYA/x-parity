package domain

import "fmt"

type ConfigValue struct {
	State EvidenceState
	Value string
	Sensitive bool
}

type SensitivePatterns string
const (
	PatternKey SensitivePatterns = "KEY"
	PatternPassword SensitivePatterns = "PASSWORD"
	PatternSecret SensitivePatterns = "SECRET"
	PatternToken SensitivePatterns = "TOKEN"
	PatternAuth SensitivePatterns = "AUTH"
	PatternCredential SensitivePatterns = "CREDENTIAL"
	PatternPrivate SensitivePatterns = "PRIVATE"
)

func (sp SensitivePatterns) IsValid() bool {
	switch sp {
	case PatternKey, PatternPassword, PatternSecret, PatternToken, PatternAuth, PatternCredential, PatternPrivate:
		return true
	default:
		return false
	}
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