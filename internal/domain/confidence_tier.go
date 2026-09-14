package domain

type ConfidenceTier int
const (
	TierDeterministic ConfidenceTier = 1
	TierDerived ConfidenceTier = 2
	TierProbabilistic ConfidenceTier = 3
	TierSpeculative ConfidenceTier = 4
)

func (ct ConfidenceTier) String() string {
	switch ct {
	case TierDeterministic:
		return "deterministic"
	case TierDerived:
		return "derived"
	case TierProbabilistic:
		return "probabilistic"
	case TierSpeculative:
		return "speculative"
	default:
		return "unknown"
	}
}

func (ct ConfidenceTier) IsValid() bool {
	switch ct {
	case TierDeterministic, TierDerived, TierProbabilistic, TierSpeculative:
		return true
	default:
		return false
	}
}