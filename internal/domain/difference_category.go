package domain

type DifferenceCategory string
const (
	CatIdentity DifferenceCategory = "identity"
	CatRuntime DifferenceCategory = "runtime"
	CatConfiguration DifferenceCategory = "configuration"
	CatSource DifferenceCategory = "source"
	CatDependency DifferenceCategory = "dependency"
)

func (c DifferenceCategory) IsValid() bool {
	switch c {
	case CatIdentity, CatRuntime, CatConfiguration, CatSource, CatDependency:
		return true
	default:
		return false
	}
}

func (c DifferenceCategory) String() string {
	return string(c)
}