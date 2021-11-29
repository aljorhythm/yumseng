package rooms

import (
	"github.com/aljorhythm/yumseng/cheers"
	"github.com/aljorhythm/yumseng/utils"
)

type CheerAddedMessage struct {
	Cheer     cheers.Cheer `json:"cheer"`
	EventName string       `json:"event_name"`
}

func NewCheerAddedMessage(cheer cheers.Cheer) ([]byte, error) {
	eventType := EVENT_CHEER_ADDED
	message := CheerAddedMessage{Cheer: cheer, EventName: string(eventType.name)}
	bytes, err := utils.ToJson(message)

	if err != nil {
		return nil, nil
	}

	return bytes, nil
}

type RoomLastSecondsCheerCountMessage struct {
	Count     int    `json:"count"`
	EventName string `json:"event_name"`
}

func NewRoomLastSecondsCheerCountMessage(count int) ([]byte, error) {
	eventType := EVENT_LAST_SECONDS_COUNT
	message := RoomLastSecondsCheerCountMessage{Count: count, EventName: string(eventType.name)}
	bytes, err := utils.ToJson(message)

	if err != nil {
		return nil, nil
	}

	return bytes, nil
}
