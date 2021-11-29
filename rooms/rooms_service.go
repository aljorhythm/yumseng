package rooms

import "github.com/aljorhythm/yumseng/cheers"

type RoomServicer interface {
	AddCheer(room *Room, cheer cheers.Cheer)
	ListenCheer(room *Room, clientId string, callback Callback)
	StopListeningCheers(room *Room, clientId string)
}

type roomsService struct {
	*RoomEvents
}

func (r *roomsService) StopListeningCheers(room *Room, clientId string) {
	r.UnsubscribeCheerAdded(room, clientId)
}

func (r *roomsService) ListenCheer(room *Room, clientId string, callback Callback) {
	r.SubscribeCheerAdded(room, clientId, callback)
}

func (r *roomsService) AddCheer(room *Room, cheer cheers.Cheer) {
	r.PublishCheerAdded(room, cheer)
}

func NewRoomsService() RoomServicer {
	return &roomsService{
		NewRoomEvents(),
	}
}
