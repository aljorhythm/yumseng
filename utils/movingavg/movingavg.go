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

type MovingAverageCalculator interface {
	AddItem(item Item)
	CountFrom(ago time.Duration) int
}

type calculator struct {
	items []Item
	Nower Nower
}

func (calculator *calculator) CountFrom(ago time.Duration) int {
	cursor := calculator.Nower.Now()
	since := cursor.Add(-ago)
	count := 0

	for i := len(calculator.items) - 1; i >= 0; i-- {
		item := calculator.items[i]
		if item.GetTime().Before(since) {
			break
		} else {
			count += 1
		}
	}

	return count
}

func (calculator *calculator) AddItem(item Item) {
	calculator.items = append(calculator.items, item)
}

type Nower interface {
	Now() time.Time
}

func NewCalculator(nower Nower) MovingAverageCalculator {
	cal := &calculator{[]Item{}, nower}
	return cal
}
