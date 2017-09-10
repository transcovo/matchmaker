package match

import "time"

var durations = []time.Duration{
	120 * time.Minute,
	90 * time.Minute,
}

var minSessionSpacing = time.Hour * 2

var maxSessionsPerWeek = 3

var maxWidthExploration = 2

var maxExplorationPathLength = 10
