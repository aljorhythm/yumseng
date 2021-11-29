package rooms

import (
	"github.com/aljorhythm/yumseng/cheers"
)

type Room struct {
	Cheers []*cheers.Cheer
	Name   string
}

func NewRoom(name string) *Room {
	return &Room{
		[]*cheers.Cheer{},
		name,
	}
}
