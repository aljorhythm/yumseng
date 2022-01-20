package movingavg

import (
	"time"
)

/**
used for calculation of average over time
usecase is calculating speed of cheers in a room
*/

type Item interface {
	GetTime() time.Time
}

type MovingAverageCalculator struct {
	Nower Nower
}

func (MovingAverageCalculator *MovingAverageCalculator) CountFrom(ago time.Duration, items []Item) int {
	cursor := MovingAverageCalculator.Nower.Now()
	since := cursor.Add(-ago)
	count := 0

	for i := len(items) - 1; i >= 0; i-- {
		item := items[i]
		if item.GetTime().Before(since) {
			break
		} else {
			count += 1
		}
	}

	return count
}

func NewCalculator(nower Nower) MovingAverageCalculator {
	cal := MovingAverageCalculator{nower}
	return cal
}
