package domain

type Severity string
const (
	SeverityCritical Severity = "critical"
	SeverityWarning Severity = "warning"
	SeverityInfo Severity = "info"
)
func (s Severity) IsValid() bool {
	switch s {
	case SeverityCritical, SeverityWarning, SeverityInfo:
		return true
	default:
		return false
	}
}

func (s Severity) String() string {
	return string(s)
}