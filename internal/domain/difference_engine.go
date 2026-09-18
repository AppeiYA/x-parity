package domain

import "strings"

type DiffEngine struct{}

func NewDiffEngine() *DiffEngine {
	return &DiffEngine{}
}

func (e *DiffEngine) Compare(local, remote Snapshot) ([]Difference, error) {
	return e.CompareWithOptions(local, remote, false)
}

func (e *DiffEngine) CompareWithOptions(local, remote Snapshot, crossPlatform bool) ([]Difference, error) {
	var diffs []Difference

	// 1. Runtime OS & Architecture
	if local.Runtime() != nil && remote.Runtime() != nil {
		if local.Runtime().OS != remote.Runtime().OS {
			sev := SeverityWarning
			msg := "Operating system platform mismatch"
			cat := CatRuntime

			if crossPlatform {
				sev = SeverityInfo
				cat = CatCrossPlatform
				msg = "Cross-platform comparison: OS divergence permitted (" + local.Runtime().OS + " vs " + remote.Runtime().OS + ")"
				if local.Runtime().Virtualization != "" || remote.Runtime().Virtualization != "" {
					msg += " [local virt: " + local.Runtime().Virtualization + ", remote virt: " + remote.Runtime().Virtualization + "]"
				}
			}

			diff, err := NewDifference(
				cat,
				"runtime.os",
				local.Runtime().OS,
				remote.Runtime().OS,
				sev,
				msg,
			)
			if err != nil {
				return nil, err
			}
			diffs = append(diffs, *diff)
		}

		// Architecture
		if local.Runtime().Architecture != remote.Runtime().Architecture {
			sev := SeverityWarning
			cat := CatRuntime
			msg := "Hardware CPU architecture mismatch"
			if crossPlatform {
				sev = SeverityInfo
				cat = CatCrossPlatform
				msg = "Cross-platform comparison: architecture divergence permitted (" + local.Runtime().Architecture + " vs " + remote.Runtime().Architecture + ")"
			}

			diff, err := NewDifference(
				cat,
				"runtime.architecture",
				local.Runtime().Architecture,
				remote.Runtime().Architecture,
				sev,
				msg,
			)
			if err != nil {
				return nil, err
			}
			diffs = append(diffs, *diff)
		}

		// Language Runtimes - strictly enforced in all modes
		for lang, localVer := range local.Runtime().Runtimes {
			remoteVer, exists := remote.Runtime().Runtimes[lang]
			if !exists {
				diff, err := NewDifference(
					CatRuntime,
					"runtime.runtimes."+lang,
					localVer,
					"missing",
					SeverityCritical,
					lang+" runtime is present locally but missing in remote",
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
					lang+" runtime version divergence detected",
				)
				if err != nil {
					return nil, err
				}
				diffs = append(diffs, *diff)
			}
		}
	}

	// 2. Configuration variables
	if local.Configuration() != nil && remote.Configuration() != nil {
		for key, localVal := range local.Configuration().Variables {
			remoteKey, remoteVal, exists := findEquivalentConfigVar(key, remote.Configuration().Variables, crossPlatform)
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
				continue
			}

			if !localVal.Sensitive && localVal.Value != remoteVal.Value {
				sev := SeverityWarning
				msg := "Non-sensitive configuration variable value mismatch"
				cat := CatConfiguration

				// In cross-platform mode, recognize host path divergence as informational
				if crossPlatform && isPathConfigKey(key, remoteKey) {
					sev = SeverityInfo
					cat = CatCrossPlatform
					msg = "Cross-platform path variable divergence (host-specific path formats permitted)"
				}

				diff, err := NewDifference(
					cat,
					"configuration.variables."+key,
					localVal.Value,
					remoteVal.Value,
					sev,
					msg,
				)
				if err != nil {
					return nil, err
				}
				diffs = append(diffs, *diff)
			}
		}
	}

	// 3. Compare source git commit
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

func findEquivalentConfigVar(key string, targetVars map[string]*ConfigValue, crossPlatform bool) (string, *ConfigValue, bool) {
	if val, exists := targetVars[key]; exists {
		return key, val, true
	}
	if !crossPlatform {
		return "", nil, false
	}

	// 1. Case-insensitive search
	lowerKey := strings.ToLower(key)
	for tKey, val := range targetVars {
		if strings.ToLower(tKey) == lowerKey {
			return tKey, val, true
		}
	}

	// 2. Canonical Semantic Cross-Platform Aliases
	aliases := map[string][]string{
		"PATH":        {"Path", "path"},
		"Path":        {"PATH", "path"},
		"HOME":        {"USERPROFILE"},
		"USERPROFILE": {"HOME"},
		"TMPDIR":      {"TEMP", "TMP"},
		"TEMP":        {"TMPDIR", "TMP"},
		"TMP":         {"TMPDIR", "TEMP"},
	}

	if targetAliases, ok := aliases[key]; ok {
		for _, alias := range targetAliases {
			if val, exists := targetVars[alias]; exists {
				return alias, val, true
			}
			lowerAlias := strings.ToLower(alias)
			for tKey, val := range targetVars {
				if strings.ToLower(tKey) == lowerAlias {
					return tKey, val, true
				}
			}
		}
	}

	return "", nil, false
}

func isPathConfigKey(k1, k2 string) bool {
	u1 := strings.ToUpper(k1)
	u2 := strings.ToUpper(k2)
	return u1 == "PATH" || u2 == "PATH" ||
		u1 == "HOME" || u2 == "HOME" ||
		u1 == "USERPROFILE" || u2 == "USERPROFILE" ||
		u1 == "TEMP" || u2 == "TEMP" ||
		u1 == "TMP" || u2 == "TMP" ||
		u1 == "TMPDIR" || u2 == "TMPDIR"
}

func isPatchVersionCompatible(v1, v2 string) bool {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	if len(parts1) >= 2 && len(parts2) >= 2 {
		return parts1[0] == parts2[0] && parts1[1] == parts2[1]
	}
	return v1 == v2
}