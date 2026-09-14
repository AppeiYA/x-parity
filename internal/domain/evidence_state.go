package domain

type EvidenceState string
const (
	StateKnown EvidenceState = "known"
	StateUnknown EvidenceState = "unknown"
	StateUnavailable EvidenceState = "unavailable"
	StateRedacted EvidenceState = "redacted"
	StateMissing EvidenceState = "missing"
)

func (e EvidenceState) IsValid() bool {
	switch e {
	case StateKnown, StateUnknown, StateUnavailable, StateRedacted, StateMissing:
		return true
	default:
		return false
	}
}

