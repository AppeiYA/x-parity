package domain

import "strings"

type DiffEngine struct{}

func NewDiffEngine() *DiffEngine {
	return &DiffEngine{}
}

func (e *DiffEngine) Compare(local, remote Snapshot) ([]Difference, error) {
	var diffs []Difference

	// Runtime OS & Architecture
	if local.Runtime() != nil && remote.Runtime() != nil {
		if local.Runtime().OS != remote.Runtime().OS {
			diff, err := NewDifference(
				CatRuntime,
				"runtime.os",
				local.Runtime().OS,
				remote.Runtime().OS,
				SeverityWarning,
				"Operating system platform mismatch",
			)
			if err != nil {
				return nil, err
			}

			diffs = append(diffs, *diff)
		}

		// Language Runtimes
		for lang, localVer := range local.Runtime().Runtimes {
			remoteVer, exists := remote.Runtime().Runtimes[lang]
			if !exists {
				diff, err := NewDifference(
					CatRuntime,
					"runtime.runtimes."+lang,
					localVer,
					"missing",
					SeverityCritical,
					lang + " runtime is present locally but missing in remote",
				)
				if err != nil {
					return nil, err
				}
				diffs = append(diffs, *diff)
				continue
			}

			if localVer != remoteVer {
				sev := SeverityWarning
				if !isPatchVersionCompatible(localVer, remoteVer) {
					sev = SeverityCritical
				}

				diff, err := NewDifference(
					CatRuntime,
					"runtime.runtimes."+lang,
					localVer,
					remoteVer,
					sev,
					lang + " runtime version divergence detected",
				)
				if err != nil {
					return nil, err
				}
				diffs = append(diffs, *diff)
			}
		}
	}

	// Configuration variables
	if local.Configuration() != nil && remote.Configuration() != nil {
		for key, localVal := range local.Configuration().Variables {
			remoteVal, exists := remote.Configuration().Variables[key]
			if !exists || remoteVal.State == StateMissing {
				diff, err := NewDifference(
					CatConfiguration,
					"configuration.variables."+key,
					localVal.State,
					StateMissing,
					SeverityCritical,
					"Configuration variable is missing in target environment",
				)
				if err != nil {
					return nil, err
				}
				diffs = append(diffs, *diff)
			}

			if !localVal.Sensitive && localVal.Value != remoteVal.Value {
				diff, err := NewDifference(
					CatConfiguration,
					"configuration.variables."+key,
					localVal.Value,
					remoteVal.Value,
					SeverityWarning,
					"Non-sensitive configuration variable value mismatch",
				)
				if err != nil {
					return nil, err
				}
				diffs = append(diffs, *diff)
			}
		}
	}

	// Compare source git commit 
	if local.Source() != nil && remote.Source() != nil {
		if local.Source().Commit != "" && remote.Source().Commit != "" && local.Source().Commit != remote.Source().Commit {
			diff, err := NewDifference(
				CatSource,
				"source.commit",
				local.Source().Commit,
				remote.Source().Commit,
				SeverityWarning,
				"Source commit mismatch between environments",
			)
			if err != nil {
				return nil, err
			}
			diffs = append(diffs, *diff)
		}
	}

	return diffs, nil
}

func isPatchVersionCompatible(v1, v2 string) bool {
	// DOMAIN RULE: Major and minor versions must match
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	if len(parts1) >= 2 && len(parts2) >= 2 {
		return parts1[0] == parts2[0] && parts1[1] == parts2[1]
	}
	return v1 == v2
}