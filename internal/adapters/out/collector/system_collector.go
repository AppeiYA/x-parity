package collector

type SystemCollector struct {}
var colectorName = "system"

func NewSystemCollector() *SystemCollector {
	return &SystemCollector{}
}

func (sc *SystemCollector) Name() string {
	return colectorName
}

