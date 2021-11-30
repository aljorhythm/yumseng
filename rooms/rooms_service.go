package rooms

import "github.com/aljorhythm/yumseng/cheers"

type RoomServicer interface {
	AddCheer(room *Room, cheer *cheers.Cheer)
	ListenCheer(room *Room, clientId string, callback Callback)
	StopListeningCheers(room *Room, clientId string)
	GetRoom(name string) *Room
}

type roomsService struct {
	*RoomEvents
	rooms map[string]*Room
}

func (r *roomsService) GetRoom(name string) *Room {
	if room, ok := r.rooms[name]; ok {
		return room
	} else {
		r.rooms[name] = NewRoom(name)
		return r.rooms[name]
	}
}

func (r *roomsService) StopListeningCheers(room *Room, clientId string) {
	r.UnsubscribeCheerAdded(room, clientId)
}

func (r *roomsService) ListenCheer(room *Room, clientId string, callback Callback) {
	r.SubscribeCheerAdded(room, clientId, callback)
}

func (r *roomsService) AddCheer(room *Room, cheer *cheers.Cheer) {
	room.AddCheer(cheer)
	r.PublishCheerAdded(room, *cheer)
}

func NewRoomsService() RoomServicer {
	return &roomsService{
		NewRoomEvents(),
		map[string]*Room{},
	}
}
