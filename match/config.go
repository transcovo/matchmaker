package match

import "time"

var durations = []time.Duration{
	60 * time.Minute,
}

var minSessionSpacing = time.Hour * 8

var maxSessionsPerWeek = 1

var maxWidthExploration = 2

var maxExplorationPathLength = 10
