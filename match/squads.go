package match

func generateSquads(people []*Person, busyTimes []*BusyTime) []*Squad {
	masters := filterPersons(people, true)
	disciples := filterPersons(people, false)

	squads := []*Squad{}
	for _, master := range masters {
		for _, disciple := range disciples {
			masterExclusivity := master.GetExclusivity()
			discipleExclusivity := disciple.GetExclusivity()
			languagesCompatible := masterExclusivity == ExclusivityNone ||
				discipleExclusivity == ExclusivityNone ||
				masterExclusivity == discipleExclusivity
			if languagesCompatible {
				people := []*Person{master, disciple}
				squads = append(squads, &Squad{
					People:     people,
					BusyRanges: mergeBusyRanges(busyTimes, people),
				})
			}
		}
	}

	return squads
}

func filterPersons(persons []*Person, wantedIsGoodReviewer bool) []*Person {
	result := []*Person{}
	for _, person := range persons {
		if person.IsGoodReviewer == wantedIsGoodReviewer {
			result = append(result, person)
		}
	}
	return result
}

func mergeBusyRanges(busyTimes []*BusyTime, people []*Person) []*Range {
	mergedBusyRanges := []*Range{}
	for _, busyTime := range busyTimes {
		for _, person := range people {
			if busyTime.Person == person {
				mergedBusyRanges = mergeRangeListWithRange(mergedBusyRanges, busyTime.Range)
			}
		}
	}
	return mergedBusyRanges
}

func mergeRangeListWithRange(ranges []*Range, extraRange *Range) []*Range {
	mergedRangeList := []*Range{}
	rangeToAdd := extraRange
	for _, currentRange := range ranges {
		if haveIntersection(currentRange, extraRange) {
			rangeToAdd = mergeRanges(currentRange, rangeToAdd)
		} else {
			mergedRangeList = append(mergedRangeList, currentRange)
		}
	}
	return append(mergedRangeList, rangeToAdd)
}

func mergeRanges(range1 *Range, range2 *Range) *Range {
	result := &Range{}
	if range1.Start.Before(range2.Start) {
		result.Start = range1.Start
	} else {
		result.Start = range2.Start
	}
	if range1.End.After(range2.End) {
		result.End = range1.End
	} else {
		result.End = range2.End
	}
	return result
}

func haveIntersection(range1 *Range, range2 *Range) bool {
	return range1.End.After(range2.Start) && range2.End.After(range1.Start)
}
