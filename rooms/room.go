package rooms

import (
	"github.com/aljorhythm/yumseng/cheers"
	"github.com/aljorhythm/yumseng/utils/movingavg"
	"time"
)

type Room struct {
	Cheers []*cheers.Cheer
	Name   string
	movingavg.MovingAverageCalculator
}

type CheerItem struct {
	cheers.Cheer
}

func (item CheerItem) GetTime() time.Time {
	return item.Cheer.ClientCreatedAt
}

func (room *Room) AddCheer(cheer *cheers.Cheer) {
	room.Cheers = append(room.Cheers, cheer)
	room.MovingAverageCalculator.AddItem(CheerItem{*cheer})
}

func NewRoom(name string) *Room {
	return &Room{
		[]*cheers.Cheer{},
		name,
		movingavg.NewCalculator(movingavg.NowTime{}),
	}
}
