package domain

import "time"

type Hypothesis struct {
	Title 	 string
	LikelyCause string
	Confidence ConfidenceTier
	Evidence  []string
	Remediation string
}

type DiagnosisReport struct {
	TargetSnapshot string 
	GeneratedAt time.Time
	Differences []Difference
	Hypotheses []Hypothesis
	Summary string
}
