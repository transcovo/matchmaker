package match

import "time"

type ReviewSession struct {
	Reviewers []Person
	Start     time.Time
	End       time.Time
}

type Solution []ReviewSession
