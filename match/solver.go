package match

import (
	"time"
	"github.com/transcovo/go-chpr-logger"
	"github.com/sirupsen/logrus"
	"strconv"
	"sort"
)

type Solution struct {
	Sessions []*ReviewSession
}

func Solve(problem *Problem) *Solution {
	squads := generateSquads(problem.People, problem.BusyTimes)
	ranges := generateTimeRanges(problem.WorkRanges)
	sessions := generateSessions(squads, ranges)

	printSquads(squads)
	printRanges(ranges)

	solution := &Solution{
		Sessions: getBestSolution(problem, sessions, []*ReviewSession{}, ""),
	}

	sort.Sort(ByStart(solution.Sessions))

	printSessions(solution.Sessions)

	return solution
}

func printSessions(sessions []*ReviewSession) {
	print(len(sessions), " session(s):")
	println()
	for _, session := range sessions {
		println(session.Reviewers.People[0].Email + "\t" +
			session.Reviewers.People[1].Email + "\t" +
			session.Range.Start.Format(time.Stamp) + "\t" +
			session.Range.End.Format(time.Stamp) + "\t")
	}
}

func getBestSolution(problem *Problem, allSessions []*ReviewSession, currentSessions []*ReviewSession, path string) []*ReviewSession {
	workRanges := problem.WorkRanges
	targetCoverage := problem.TargetCoverage

	bestMissingCoverage := getMissingCoverage(currentSessions, workRanges, targetCoverage)
	bestSolutionSessions := currentSessions
	bestPath := ""

	for i, session := range allSessions {
		subpath := path + "/" + strconv.Itoa(i)

		sessionCompatible := isSessionCompatible(session, currentSessions)
		l := logger.WithField("path", subpath)

		if sessionCompatible {

			newSessions := append(currentSessions, session)
			newMissingCoverage := getMissingCoverage(newSessions, workRanges, targetCoverage)

			if isEnough(newMissingCoverage) {
				l.Info("Missing coverage ok returning solution")
				return newSessions
			}

			coverageImproved := isMissingCoverageBetter2(newMissingCoverage, bestMissingCoverage)
			l.WithFields(logrus.Fields{
				"coverageImproved": coverageImproved,
				"best":             missingCoverageToString(bestMissingCoverage),
				"new":              missingCoverageToString(newMissingCoverage),
			}).Info("Coverage comparision")

			if coverageImproved {
				bestMissingCoverage = newMissingCoverage
				bestSolutionSessions = newSessions
				bestPath = subpath
			}
		}
	}

	if bestPath != "" {
		return getBestSolution(problem, allSessions, bestSolutionSessions, bestPath)
	}

	return bestSolutionSessions
}

func missingCoverageToString(missingCoverage map[Exclusivity]int) string {
	return "[" + strconv.Itoa(missingCoverage[ExclusivityMobile]) + "/" +
		strconv.Itoa(missingCoverage[ExclusivityBack]) + "//" +
		strconv.Itoa(missingCoverage[ExclusivityNone]) + "]"
}

func isMissingCoverageBetter2(coverage1 map[Exclusivity]int, coverage2 map[Exclusivity]int) bool {
	return coverage1[ExclusivityNone] < coverage2[ExclusivityNone]
}

func isMissingCoverageBetter(coverage1 map[Exclusivity]int, coverage2 map[Exclusivity]int) bool {
	return coverage1[ExclusivityNone] <= coverage2[ExclusivityNone] &&
		coverage1[ExclusivityBack] <= coverage2[ExclusivityBack] &&
		coverage1[ExclusivityMobile] <= coverage2[ExclusivityMobile] && (
		coverage1[ExclusivityNone] < coverage2[ExclusivityNone] ||
			coverage1[ExclusivityBack] < coverage2[ExclusivityBack] ||
			coverage1[ExclusivityMobile] < coverage2[ExclusivityMobile])
}

func isEnough(missingCoverage map[Exclusivity]int) bool {
	for _, missing := range missingCoverage {
		if missing > 0 {
			return false
		}
	}
	return true
}

func isSessionCompatible(session *ReviewSession, sessions []*ReviewSession) bool {
	personSessionCount := map[string]int{}

	for _, otherSession := range sessions {
		if session == otherSession {
			return false
		}
		if session.Reviewers == otherSession.Reviewers {
			return false
		}
		if haveIntersection(session.Range, otherSession.Range) {
			return false
		}
		for _, person := range otherSession.Reviewers.People {
			if _, ok := personSessionCount[person.Email]; !ok {
				personSessionCount[person.Email] = 0
			}
			personSessionCount[person.Email]++
		}
	}
	for _, person := range session.Reviewers.People {
		if count, ok := personSessionCount[person.Email]; ok {
			if count >= 3 {
				return false
			}
		}
	}

	return true
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

func getMissingCoverage(sessions []*ReviewSession, workRanges []*Range, target map[Exclusivity]int) map[Exclusivity]int {
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

	for _, session := range sessions {
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
