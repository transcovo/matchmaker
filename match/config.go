package match

import "time"

var durations = []time.Duration{
	90 * time.Minute,
	60 * time.Minute,
}

var minSessionSpacing = time.Hour * 4

var maxSessionsPerWeek = 3

var maxWidthExploration = 2

var maxExplorationPathLength = 10
