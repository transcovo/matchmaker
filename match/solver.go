package match

import "time"

type Squad struct {
	People     []*Person
	BusyRanges []*Range
}

func (squad *Squad) GetExclusivity() Exclusivity {
	for _, person := range squad.People {
		exclusivity := person.GetExclusivity()
		if exclusivity != ExclusivityNone {
			return exclusivity
		}
	}
	return ExclusivityNone
}

type ReviewSession struct {
	Reviewers *Squad
	Start     time.Time
	End       time.Time
}

type Score struct {
	Hours    int
	Coverage float32
}

type Solution []ReviewSession

func Solve(problem *Problem) *Solution {
	squads := generateSquads(problem.People, problem.BusyTimes)
	printSquads(squads)
	return nil
}

func printSquads(squads []*Squad) {
	for _, squad := range squads {
		exclusivityStr := squad.GetExclusivity().String()
		println(squad.People[0].Email + " + " + squad.People[1].Email + ": " + exclusivityStr)
	}
}
