package domain

import "fmt"

type Difference struct {
	category DifferenceCategory
	path string
	local any
	remote any
	severity Severity
	message string
}

var (
	ErrInvalidDifferenceCategory = fmt.Errorf("invalid difference category")
	ErrInvalidSeverity = fmt.Errorf("invalid severity")
)

func NewDifference(category DifferenceCategory, path string, local, remote any, severity Severity, message string) (*Difference, error) {
	if !category.IsValid() {
		return nil, ErrInvalidDifferenceCategory
	}
	if !severity.IsValid() {
		return nil, ErrInvalidSeverity
	}
	return &Difference{
		category: category,
		path: path,
		local: local,
		remote: remote,
		severity: severity,
		message: message,
	}, nil
}

func (d *Difference) Category() DifferenceCategory {
	return d.category
}

func (d *Difference) Path() string {
	return d.path
}

func (d *Difference) Local() any {
	return d.local
}

func (d *Difference) Remote() any {
	return d.remote
}

func (d *Difference) Severity() Severity {
	return d.severity
}

func (d *Difference) Message() string {
	return d.message
}