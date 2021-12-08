package rooms

import (
	"github.com/aljorhythm/yumseng/cheers"
)

/**
Message definitions with client
*/

type RoomConnectedMessage struct {
	EventName string `json:"event_name"`
	UserId    string `json:"user_id"`
	RoomName  string `json:"room_name"`
}

type CheerAddedMessage struct {
	Cheer     cheers.Cheer `json:"cheer"`
	EventName string       `json:"event_name"`
}

type RoomLastSecondsCheerCountMessage struct {
	Count     int    `json:"count"`
	EventName string `json:"event_name"`
}
