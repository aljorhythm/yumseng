package rooms

import (
	"context"
	"fmt"
	"github.com/aljorhythm/yumseng/cheers"
	"github.com/aljorhythm/yumseng/objectstorage"
	"log"
)

type CheerImage struct {
	Url      string `json:"url"`
	ObjectId string `json:"object-id"`
}

type RoomServicer interface {
	AddCheerImage(ctx context.Context, roomId string, user User, data []byte, id string) error
	GetCheerImages(ctx context.Context, roomId string, user User) ([]*CheerImage, error)
	UserJoinsRoom(ctx context.Context, room *Room, user User) error
	AddCheer(room *Room, cheer *cheers.Cheer, user User)
	AddCheerAddedListener(room *Room, user User, clientId string, callback Callback) error
	StopListeningCheers(room *Room, clientId string)
	GetOrCreateRoom(name string) *Room
	GetRoom(name string) *Room
}

type roomsService struct {
	*RoomEvents
	rooms         map[string]*Room
	objectStorage objectstorage.Storage
}

func (r *roomsService) GetRoom(name string) *Room {
	return r.rooms[name]
}

func (r *roomsService) UserJoinsRoom(ctx context.Context, room *Room, user User) error {
	_, err := room.AddUserIfNotPresent(user)
	return err
}

func cheerObjectId(room *Room, user User, id string) string {
	return fmt.Sprintf("rooms/%s/%s/%s", room.Name, user.GetId(), id)
}

func (r *roomsService) GetCheerImages(ctx context.Context, roomId string, user User) ([]*CheerImage, error) {
	room := r.GetRoom(roomId)

	if room == nil {
		return nil, ERROR_ROOM_NOT_FOUND
	}
	return room.GetCheerImages(user)
}

func (r *roomsService) AddCheerImage(ctx context.Context, roomId string, user User, data []byte, id string) error {
	room := r.GetRoom(roomId)

	if room == nil {
		return ERROR_ROOM_NOT_FOUND
	}

	objectId := cheerObjectId(room, user, id)

	if err := r.objectStorage.Store(ctx, objectId, data); err != nil {
		return err
	}

	meta, err := r.objectStorage.RetrieveObjectMetadata(ctx, objectId)
	if err != nil {
		return err
	}

	cheerImage := CheerImage{
		Url:      meta.Url,
		ObjectId: objectId,
	}

	return room.AddCheerImage(user, &cheerImage)
}

func (r *roomsService) GetOrCreateRoom(name string) *Room {
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

func (r *roomsService) AddCheerAddedListener(room *Room, user User, clientId string, callback Callback) error {
	if created, err := room.AddUserIfNotPresent(user); err != nil {
		return err
	} else {
		userId := user.GetId()
		var userQueryStatus string
		if created == true {
			userQueryStatus = fmt.Sprintf("User %s added to room %s", userId, room.Name)
		} else {
			userQueryStatus = fmt.Sprintf("User %s found in room %s", userId, room.Name)
		}

		log.Printf("EventsSocketId: %s room: %s , user id: %s, , %s ", clientId, room.Name, user.GetId(), userQueryStatus)
	}
	r.RoomEvents.SubscribeCheerAdded(room, clientId, callback)
	log.Printf("EventsSocketId[%s] userId : %s in room %s (Subscribed to cheers)", clientId, user.GetId(), room.Name)
	return nil
}

func (r *roomsService) AddCheer(room *Room, cheer *cheers.Cheer, user User) {
	cheer.UserId = user.GetId()
	room.AddCheer(cheer)
	r.RoomEvents.PublishCheerAdded(room, *cheer)
}

func NewRoomsService(storage objectstorage.Storage) RoomServicer {
	return &roomsService{
		NewRoomEvents(),
		map[string]*Room{},
		storage,
	}
}
