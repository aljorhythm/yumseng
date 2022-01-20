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

func NewRoomConnectedMessage(room *Room, user User) (*RoomConnectedMessage, error) {
	eventType := EVENT_ROOM_CONNECTED
	message := RoomConnectedMessage{RoomName: room.Name, UserId: user.GetId(), EventName: eventType.GetName()}
	return &message, nil
}

type CheerAddedMessage struct {
	Cheer     cheers.Cheer `json:"cheer"`
	EventName string       `json:"event_name"`
}

func NewCheerAddedMessage(cheer cheers.Cheer) (*CheerAddedMessage, error) {
	eventType := EVENT_CHEER_ADDED
	message := CheerAddedMessage{Cheer: cheer, EventName: eventType.GetName()}
	return &message, nil
}

type RoomIntensityMessage struct {
	Intensity int    `json:"count"`
	EventName string `json:"event_name"`
}

func NewRoomIntensityMessage(count int) (*RoomIntensityMessage, error) {
	eventType := EVENT_ROOM_INTENSITY
	message := RoomIntensityMessage{Intensity: count, EventName: eventType.GetName()}
	return &message, nil
}
