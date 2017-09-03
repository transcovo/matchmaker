package match

import (
	"time"
	"strconv"
	"sort"
	"github.com/transcovo/go-chpr-logger"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strings"
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

	bestSessions, _ := getSolver(problem, sessions)([]*ReviewSession{}, "")
	solution := &Solution{
		Sessions: bestSessions,
	}

	coverage, maxCoverage := getCoverage(problem.WorkRanges, bestSessions)

	missingCoverage := getMissingConverage(coverage, problem.TargetCoverage)

	worstMissingCoverage, _ := getCoveragePerformance([]*ReviewSession{}, problem.WorkRanges, problem.TargetCoverage)

	println(missingCoverageToString(missingCoverage))
	println(missingCoverageToString(worstMissingCoverage))

	println(maxCoverage)

	sort.Sort(ByStart(solution.Sessions))

	printSessions(solution.Sessions)

	for exclusivity, exclusivityCoverage := range coverage {
		println(exclusivity.String() + ":")
		for i, value := range exclusivityCoverage {
			println("  " + strconv.Itoa(i) + " -> " + strconv.Itoa(value))
		}
	}

	return solution
}

func printSessions(sessions []*ReviewSession) {
	print(len(sessions), " session(s):")
	println()
	for _, session := range sessions {
		printSession(session)
	}
}
func printSession(session *ReviewSession) {
	println(session.Reviewers.People[0].Email + "\t" +
		session.Reviewers.People[1].Email + "\t" +
		session.Range.Start.Format(time.Stamp) + "\t" +
		session.Range.End.Format(time.Stamp) + "\t")
}

type solver func([]*ReviewSession, string) ([]*ReviewSession, map[Exclusivity]int)

func getSolver(problem *Problem, allSessions []*ReviewSession) solver {
	var solve solver

	workRanges := problem.WorkRanges
	targetCoverage := problem.TargetCoverage

	bestSessions := []*ReviewSession{}
	bestCoveragePerformance, _ := getCoveragePerformance(bestSessions, workRanges, targetCoverage)

	var iterations int64 = 0

	solve = func(currentSessions []*ReviewSession, path string) ([]*ReviewSession, map[Exclusivity]int) {
		currentCoveragePerformance, _ := getCoveragePerformance(currentSessions, workRanges, targetCoverage)

		//for i, session := range allSessions {
		maxIter := 5000 / (len(currentSessions) + 1)
		for k := 0; k < maxIter; k++ {
			if iterations > 10e7 {
				break
			}

			i := rand.Intn(len(allSessions))
			session := allSessions[i]
			iterations += 1

			subPath := path + "/" + strconv.Itoa(i)
			if iterations%10e4 == 0 {
				logger.WithFields(logrus.Fields{
					"iterations": iterations / 1000,
					"best":       missingCoverageToString(bestCoveragePerformance),
					"current":    missingCoverageToString(currentCoveragePerformance),
				}).Info("Coverage comparision")
			}

			sessionCompatible := isSessionCompatible(session, currentSessions)
			if !sessionCompatible {
				continue
			}

			newSessions := append(currentSessions, session)
			newCoveragePerformance, newMaxCoverage := getCoveragePerformance(newSessions, workRanges, targetCoverage)
			if newMaxCoverage > problem.MaxTotalCoverage {
				continue
			}

			improvesCoverage := isMissingCoverageBetter(newCoveragePerformance, currentCoveragePerformance)
			if !improvesCoverage {
				continue
			}

			currentCoveragePerformance = newCoveragePerformance

			if isMissingCoverageBetter(newCoveragePerformance, bestCoveragePerformance) {
				bestSessions = make([]*ReviewSession, len(newSessions))
				copy(bestSessions, newSessions)
				bestCoveragePerformance = newCoveragePerformance
			}

			if isEnough(currentCoveragePerformance) {
				break
			}

			solve(newSessions, subPath)
		}

		__debug_perf, _ := getCoveragePerformance(bestSessions, workRanges, targetCoverage)
		if missingCoverageToString(__debug_perf) != missingCoverageToString(bestCoveragePerformance) {
			printSessions(bestSessions)
			println(missingCoverageToString(__debug_perf))
			println(missingCoverageToString(bestCoveragePerformance))
			panic("here")
		}

		return bestSessions, bestCoveragePerformance
	}
	return solve;
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
	// store the number of sessions to cap it
	reviewers := session.Reviewers
	people := reviewers.People

	person0 := people[0]
	person0.isSessionCompatibleSessionCount = 0
	person1 := people[1]
	person1.isSessionCompatibleSessionCount = 0

	for _, otherSession := range sessions {
		// not the same session two times
		if session == otherSession {
			return false
		}

		otherReviewers := otherSession.Reviewers
		// not the same squad
		if reviewers == otherReviewers {
			return false
		}

		otherPeople := otherReviewers.People

		otherPerson0 := otherPeople[0]
		otherPerson1 := otherPeople[1]
		if (otherPerson0 == person0 || otherPerson0 == person1 || otherPerson1 == person0 || otherPerson1 == person1) &&
			haveIntersection(session.Range, otherSession.Range) {

			return false
		}
		// every reviewer must be able to attempt all the sessions
		otherPerson0.isSessionCompatibleSessionCount += 1
		otherPerson1.isSessionCompatibleSessionCount += 1
	}

	// max 4 reviews per person
	return person0.isSessionCompatibleSessionCount < 4 && person1.isSessionCompatibleSessionCount < 4
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

func (squad *Squad) GetDisplayName() string {
	result := ""
	for _, person := range squad.People {
		if result != "" {
			result = result + " / "
		}
		result = result + getNameFromEmail(person.Email)
	}
	return result
}

func getNameFromEmail(email string) string {
	beforeA := strings.Split(email, "@")[0]
	firstName := strings.Split(beforeA, ".")[0]
	return strings.Title(firstName)
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

func getCoveragePerformance(sessions []*ReviewSession, workRanges []*Range, target map[Exclusivity]int) (map[Exclusivity]int, int) {
	coverage, maxCoverage := getCoverage(workRanges, sessions)

	missingCoverage := getMissingConverage(coverage, target)

	return missingCoverage, maxCoverage
}
func getMissingConverage(coverage map[Exclusivity]map[int]int, target map[Exclusivity]int) map[Exclusivity]int {
	missingCoverage := map[Exclusivity]int{
		ExclusivityMobile: 0,
		ExclusivityBack:   0,
		ExclusivityNone:   0,
	}
	for exclusivity, exclusivityCoverage := range coverage {
		targetValue := target[exclusivity]
		for _, value := range exclusivityCoverage {
			if value < targetValue {
				missingCoverage[exclusivity] += targetValue - value
			}
		}
	}
	return missingCoverage
}

func getCoverage(workRanges []*Range, sessions []*ReviewSession) (map[Exclusivity]map[int]int, int) {
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
	maxCoverage := 0
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
			if coverage[ExclusivityNone][coveragePeriodId] > maxCoverage {
				maxCoverage = coverage[ExclusivityNone][coveragePeriodId]
			}
			date = date.Add(coveragePeriodSpan)
		}
	}
	return coverage, maxCoverage
}

func getCoveragePeriodId(workRanges []*Range, date time.Time) int {
	elapsedNanoseconds := date.Sub(workRanges[0].Start).Nanoseconds()
	elapsedCoveragePeriods := elapsedNanoseconds / (30 * 60 * 1000 * 1000 * 1000)
	return int(elapsedCoveragePeriods)
}
