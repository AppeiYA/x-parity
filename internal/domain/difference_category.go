package domain

type DifferenceCategory string

const (
	CatIdentity       DifferenceCategory = "identity"
	CatRuntime        DifferenceCategory = "runtime"
	CatConfiguration  DifferenceCategory = "configuration"
	CatSource         DifferenceCategory = "source"
	CatDependency     DifferenceCategory = "dependency"
	CatCrossPlatform  DifferenceCategory = "cross_platform"
)

func (c DifferenceCategory) IsValid() bool {
	switch c {
	case CatIdentity, CatRuntime, CatConfiguration, CatSource, CatDependency, CatCrossPlatform:
		return true
	default:
		return false
	}
}

func (c DifferenceCategory) String() string {
	return string(c)
}