package match

import (
	"time"
)

type Solution struct {
	Sessions []ReviewSession
}

func Solve(problem *Problem) *Solution {
	squads := generateSquads(problem.People, problem.BusyTimes)
	ranges := generateTimeRanges(problem.WorkRanges)
	//sessions := generateSessions(squads, ranges)

	printSquads(squads)
	printRanges(ranges)
/*
	solution := &Solution{
		Sessions:[]ReviewSession{},
	}

	ImproveSolution(solution, sessions)*/

	return nil
}

func printRanges(ranges []*Range) {
	for _, currentRange := range ranges {
		println(currentRange.Start.Format(time.RFC3339) + " -> " + currentRange.End.Format(time.RFC3339))
	}
}

func printSquads(squads []*Squad) {
	for _, squad := range squads {
		exclusivityStr := squad.GetExclusivity().String()
		println(squad.People[0].Email + " + " + squad.People[1].Email + ": " + exclusivityStr)
	}
}

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

type Score struct {
	Hours    int
	Coverage float32
}

var coveragePeriodSpan = 30 * time.Minute

func (solution *Solution) GetMissingCoverage(workRanges []*Range, target map[Exclusivity]int) map[Exclusivity]int {
	coverage := map[Exclusivity]map[int]int{
		ExclusivityMobile: {},
		ExclusivityBack:   {},
		ExclusivityNone:   {},
	}

	for _, workRange := range workRanges {
		date := workRange.Start
		for date.Before(workRange.End) {
			coveragePeriodId := getCoveragePeriodId(workRanges, date)
			coverage[ExclusivityMobile][coveragePeriodId] = 0
			coverage[ExclusivityBack][coveragePeriodId] = 0
			coverage[ExclusivityNone][coveragePeriodId] = 0
			date = date.Add(coveragePeriodSpan)
		}
	}

	for _, session := range solution.Sessions {
		date := session.Start()
		for date.Before(session.End()) {
			coveragePeriodId := getCoveragePeriodId(workRanges, date)
			switch session.Reviewers.GetExclusivity() {
			case ExclusivityMobile:
				coverage[ExclusivityMobile][coveragePeriodId] += 1
			case ExclusivityBack:
				coverage[ExclusivityBack][coveragePeriodId] += 1
			case ExclusivityNone:
				coverage[ExclusivityMobile][coveragePeriodId] += 1
				coverage[ExclusivityBack][coveragePeriodId] += 1
			}
			coverage[ExclusivityNone][coveragePeriodId] += 1
			date = date.Add(coveragePeriodSpan)
		}
	}

	missingCoverage := map[Exclusivity]int{
		ExclusivityMobile: 0,
		ExclusivityBack:   0,
		ExclusivityNone:   0,
	}

	for exclusivity, exclusivityCoverage := range coverage {
		targetValue := target[exclusivity]
		for _, value := range exclusivityCoverage {
			if value < targetValue {
				missingCoverage[exclusivity] += (targetValue - value)
			}
		}
	}

	return missingCoverage
}

func getCoveragePeriodId(workRanges []*Range, date time.Time) int {
	elapsedNanoseconds := date.Sub(workRanges[0].Start).Nanoseconds()
	elapsedCoveragePeriods := elapsedNanoseconds / (30 * 60 * 1000 * 1000 * 1000)
	return int(elapsedCoveragePeriods)
}
