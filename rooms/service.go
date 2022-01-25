package rooms

import (
	"context"
	"fmt"
	"github.com/aljorhythm/yumseng/cheers"
	"github.com/aljorhythm/yumseng/objectstorage"
	"github.com/google/uuid"
	"log"
)

type CheerImage struct {
	Url      string `json:"url"`
	ObjectId string `json:"object_id"`
}

// mockgen -source=service.go -destination service_mockgen.go -package rooms
type RoomServicer interface {
	AddCheerImage(ctx context.Context, roomId string, user User, url string) error
	UploadCheerImage(ctx context.Context, roomId string, user User, data []byte) (*CheerImage, error)
	GetCheerImages(ctx context.Context, roomId string, user User) ([]*CheerImage, error)
	UserJoinsRoom(ctx context.Context, room *Room, user User) error
	AddCheer(room *Room, cheer *cheers.Cheer, user User) error
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

func (r *roomsService) AddCheerImage(ctx context.Context, roomId string, user User, url string) error {
	cheerImage := CheerImage{
		Url:      url,
		ObjectId: url,
	}

	room := r.GetRoom(roomId)

	err := room.AddCheerImage(user, &cheerImage)

	if err != nil {
		return err
	}

	return nil
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

func (r *roomsService) UploadCheerImage(ctx context.Context, roomId string, user User, data []byte) (*CheerImage, error) {
	room := r.GetRoom(roomId)

	if room == nil {
		return nil, ERROR_ROOM_NOT_FOUND
	}

	id := uuid.New().String()

	objectId := cheerObjectId(room, user, id)

	if err := r.objectStorage.Store(ctx, objectId, data); err != nil {
		return nil, err
	}

	meta, err := r.objectStorage.RetrieveObjectMetadata(ctx, objectId)
	if err != nil {
		return nil, err
	}

	cheerImage := CheerImage{
		Url:      meta.Url,
		ObjectId: objectId,
	}

	err = room.AddCheerImage(user, &cheerImage)

	if err != nil {
		return nil, err
	}

	return &cheerImage, nil
}

func (r *roomsService) GetOrCreateRoom(name string) *Room {
	if room, ok := r.rooms[name]; ok {
		return room
	} else {
		newRoom := NewRoom(name)
		r.rooms[name] = newRoom
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

func (r *roomsService) AddCheer(room *Room, cheer *cheers.Cheer, user User) error {
	cheer.UserId = user.GetId()
	if err := room.AddCheer(cheer); err != nil {
		return err
	}
	r.RoomEvents.PublishCheerAdded(room, *cheer)
	return nil
}

func NewRoomsService(storage objectstorage.Storage) RoomServicer {
	return &roomsService{
		NewRoomEvents(),
		map[string]*Room{},
		storage,
	}
}
