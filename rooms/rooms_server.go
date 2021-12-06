package rooms

import (
	"bytes"
	"fmt"
	"github.com/aljorhythm/yumseng/cheers"
	"github.com/aljorhythm/yumseng/objectstorage"
	"github.com/aljorhythm/yumseng/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type RoomsServer struct {
	http.Handler
	RoomServicer
	UserService UserServicer
}

type JoinRoomRequest struct {
	UserId   string `json:"user_id"`
	RoomName string `json:"room_name"`
}

func (roomsServer *RoomsServer) eventsWs(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		//todo security debt
		CheckOrigin: func(r *http.Request) bool {

			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Panicf("webSocket upgrade error %#v", err)
	}

	_, msg, err := conn.ReadMessage()
	joinRoomRequest := JoinRoomRequest{}
	err = utils.DecodeJsonFromBytes(msg, &joinRoomRequest)
	if err != nil {
		log.Panicf("unable to join room %s", string(msg))
	}

	roomName := joinRoomRequest.RoomName
	user, _ := roomsServer.UserService.GetUser(joinRoomRequest.UserId)
	room := roomsServer.RoomServicer.GetOrCreateRoom(roomName)

	clientId := fmt.Sprintf("user=%s uuid=%s", user.GetId(), uuid.New().String())

	if err != nil {
		log.Panicf("error emitting room connected %s %#v", clientId, err)
	}

	go func() {
		log.Printf("client %s listening for room %s cheers", clientId, room.Name)
		for {
			_, msg, err := conn.ReadMessage()

			if err != nil {
				log.Printf("error in reading socket message room: %s client: %s err: %#v", room.Name, clientId, err)
				return
			}
			reader := bytes.NewReader(msg)
			newCheer := cheers.Cheer{}
			utils.DecodeJson(reader, &newCheer)
			log.Printf("adding cheer from client %#v", newCheer)
			roomsServer.RoomServicer.AddCheer(room, &newCheer, user)
		}
	}()

	addedCheersChannel := make(chan cheers.Cheer)

	log.Printf("subscribing user %s client %s to room %s cheers", user.GetId(), clientId, room.Name)
	cb := func(args ...interface{}) {
		rawCheer := args[0]
		cheer, ok := rawCheer.(cheers.Cheer)
		if ok {
			log.Printf("cheer listened %#v", cheer)
			addedCheersChannel <- cheer
		} else {
			log.Panicf("cannot convert cheer %#v", args)
		}
	}
	err = roomsServer.ListenCheer(room, user, clientId, cb)

	if err != nil {
		log.Panicf("unable to subscribe user %s to room %s error %#v", user.GetId(), room.Name, err)
	}

	log.Printf("joined room %s | clientId : %s | userId : %s", room.Name, clientId, user.GetId())
	roomConnectedMessage, _ := NewRoomConnectedMessage(room, user)
	err = conn.WriteJSON(roomConnectedMessage)

	log.Printf("client %s listening to room %s cheer speed", clientId, room.Name)
	ticker := time.NewTicker(250 * time.Millisecond)
	quitIntensityListener := make(chan struct{})

	conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("exiting listening client %s code %d text %s", clientId, code, text)
		roomsServer.StopListeningCheers(room, clientId)
		close(quitIntensityListener)
		return nil
	})

	for {
		select {
		case cheer, more := <-addedCheersChannel:
			if more {
				cheerAddedMessage, err := NewCheerAddedMessage(cheer)
				log.Printf("%s writing to socket %#v", clientId, cheer)
				err = conn.WriteJSON(cheerAddedMessage)
				if err != nil {
					log.Panicf("client %s webSocket erroring write message %#v", clientId, err)
				}
			} else {
				log.Printf("cheers channel is closed %s", clientId)
			}
		case <-ticker.C:
			count := room.CountFrom((time.Duration(1) * time.Second))
			message, _ := NewRoomLastSecondsCheerCountMessage(count)
			err = conn.WriteJSON(message)
			if err != nil {
				log.Printf("err writing to socket %#v closing quit channel %s", err, clientId)
			} else {
				log.Printf("wrote to socket last seconds cheer count %s %d", clientId, count)
			}
		case <-quitIntensityListener:
			log.Printf("quit channel emitted stopping speed ticker %s", clientId)
			ticker.Stop()
			return
		}
	}
}

func NewRoomsServer(router *mux.Router, userService UserServicer, storage objectstorage.Storage) http.Handler {
	roomsServer := &RoomsServer{
		RoomServicer: NewRoomsService(storage),
		UserService:  userService,
	}

	router.Handle("/rooms", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "you are at rooms")
	}))

	router.Handle("/events", http.HandlerFunc(roomsServer.eventsWs))

	roomsServer.Handler = router
	return roomsServer
}
