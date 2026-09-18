package collector

import (
	"context"
	"os"
	"strings"

	"github.com/AppeiYA/x-parity/internal/domain"
	portout "github.com/AppeiYA/x-parity/internal/ports/out"
)

var _ portout.CollectorInt = (*ConfigCollector)(nil)

type ConfigCollector struct{}

func NewConfigCollector() *ConfigCollector {
	return &ConfigCollector{}
}

func (cc *ConfigCollector) Name() string {
	return portout.CollectorConfig
}

func (cc *ConfigCollector) Collect(ctx context.Context, s *domain.Snapshot) error {
	cfg := s.Configuration()
	if cfg == nil {
		var err error
		cfg, err = domain.NewConfiguration(make(map[string]*domain.ConfigValue))
		if err != nil {
			return err
		}
		s.SetConfiguration(cfg)
	}
	if cfg.Variables == nil {
		cfg.Variables = make(map[string]*domain.ConfigValue)
	}

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key, val := parts[0], parts[1]
		val = strings.TrimRight(val, "\r")
		isSecret := isSensitiveKey(key)
		if isSecret {
			cfg.Variables[key] = &domain.ConfigValue{
				State:     domain.StateRedacted,
				Value:     "[REDACTED]",
				Sensitive: true,
			}
		} else {
			cfg.Variables[key] = &domain.ConfigValue{
				State:     domain.StateKnown,
				Value:     val,
				Sensitive: false,
			}
		}
	}

	return nil
}

func isSensitiveKey(key string) bool {
	upper := strings.ToUpper(key)
	patterns := []string{
		string(domain.PatternKey),
		string(domain.PatternPassword),
		string(domain.PatternSecret),
		string(domain.PatternToken),
		string(domain.PatternAuth),
		string(domain.PatternCredential),
		string(domain.PatternPrivate),
	}
	for _, p := range patterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}