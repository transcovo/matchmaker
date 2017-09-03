package match

import (
	"time"
)

type ReviewSession struct {
	Reviewers *Squad
	Range     *Range
}

func (session *ReviewSession) End() time.Time {
	return session.Range.End
}

func (session *ReviewSession) Start() time.Time {
	return session.Range.Start
}

func generateSessions(squads []*Squad, ranges []*Range) []*ReviewSession {
	sessions := []*ReviewSession{}
	for _, currentRange := range ranges {
		for _, squad := range squads {
				sessions = append(sessions, &ReviewSession{
					Reviewers:squad,
					Range:currentRange,
				})
		}
	}
	return sessions
}
