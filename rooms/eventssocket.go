package rooms

import (
	"bytes"
	"fmt"
	"github.com/aljorhythm/yumseng/cheers"
	"github.com/aljorhythm/yumseng/utils"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type eventsSocket struct {
	conn                  *websocket.Conn
	roomsServer           *RoomsServer
	user                  User
	room                  *Room
	clientId              string
	addedCheersChannel    chan cheers.Cheer
	quitIntensityListener chan (struct{})
}

type JoinRoomRequest struct {
	UserId   string `json:"user_id"`
	RoomName string `json:"room_name"`
}

func (s *eventsSocket) processFirstMessage() {
	_, msg, err := s.conn.ReadMessage()
	joinRoomRequest := JoinRoomRequest{}
	err = utils.DecodeJsonFromBytes(msg, &joinRoomRequest)
	if err != nil {
		log.Panicf("unable to join room %s", string(msg))
	}

	roomName := joinRoomRequest.RoomName
	user, _ := s.roomsServer.UserService.GetUser(joinRoomRequest.UserId)
	s.user = user
	room := s.roomsServer.RoomServicer.GetOrCreateRoom(roomName)
	s.room = room
	clientId := fmt.Sprintf("user=%s uuid=%s", user.GetId(), uuid.New().String())
	s.clientId = clientId

	if err != nil {
		log.Panicf("error emitting room connected %s %#v", clientId, err)
	}
}

func (s *eventsSocket) listenToClientMessages() {
	go func() {

		log.Printf("client %s listening for room %s cheers", s.clientId, s.room.Name)
		for {
			_, msg, err := s.conn.ReadMessage()

			if err != nil {
				log.Printf("error in reading socket message room: %s client: %s err: %#v", s.room.Name, s.clientId, err)
				return
			}
			reader := bytes.NewReader(msg)
			newCheer := cheers.Cheer{}
			utils.DecodeJson(reader, &newCheer)
			log.Printf("adding cheer from client %#v", newCheer)
			s.roomsServer.RoomServicer.AddCheer(s.room, &newCheer, s.user)
		}
	}()
}

func (socket *eventsSocket) listenToRoomCheers() {
	log.Printf("subscribing user %s client %s to room %s cheers", socket.user.GetId(), socket.clientId, socket.room.Name)
	callback := func(args ...interface{}) {
		rawCheer := args[0]
		cheer, ok := rawCheer.(cheers.Cheer)
		if ok {
			log.Printf("cheer listened %#v", cheer)
			socket.addedCheersChannel <- cheer
		} else {
			log.Panicf("cannot convert cheer %#v", args)
		}
	}

	var err error
	err = socket.roomsServer.AddCheerAddedListener(socket.room, socket.user, socket.clientId, callback)

	if err != nil {
		log.Panicf("unable to subscribe user %s to room %s error %#v", socket.user.GetId(), socket.room.Name, err)
	}
}

func (socket *eventsSocket) sendRoomConnnectedMessage() {
	roomConnectedMessage, _ := NewRoomConnectedMessage(socket.room, socket.user)
	err := socket.conn.WriteJSON(roomConnectedMessage)

	if err != nil {
		log.Panicf("%#v %#v", socket, err)
	}
}

func (socket *eventsSocket) handleEventsAndSendMessages() {
	ticker := time.NewTicker(250 * time.Millisecond)
	for {
		select {
		case cheer, more := <-socket.addedCheersChannel:
			if more {
				cheerAddedMessage, err := NewCheerAddedMessage(cheer)
				log.Printf("%s writing to socket %#v", socket.clientId, cheer)
				err = socket.conn.WriteJSON(cheerAddedMessage)
				if err != nil {
					log.Panicf("client %s webSocket erroring write message %#v", socket.clientId, err)
				}
			} else {
				log.Printf("cheers channel is closed %s", socket.clientId)
			}
		case <-ticker.C:
			count := socket.room.CountFrom((time.Duration(1) * time.Second))
			message, _ := NewRoomLastSecondsCheerCountMessage(count)
			err := socket.conn.WriteJSON(message)
			if err != nil {
				log.Printf("err writing to socket %#v closing quit channel %s", err, socket.clientId)
			} else {
				log.Printf("wrote to socket last seconds cheer count %s %d", socket.clientId, count)
			}
		case <-socket.quitIntensityListener:
			log.Printf("quit channel emitted stopping speed ticker %s", socket.clientId)
			ticker.Stop()
			return
		}
	}

}

func (s *eventsSocket) setCloseHandler() {
	s.conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("exiting listening client %s code %d text %s", s.clientId, code, text)
		s.roomsServer.StopListeningCheers(s.room, s.clientId)
		close(s.quitIntensityListener)
		return nil
	})
}

func NewRoomConnectedMessage(room *Room, user User) (*RoomConnectedMessage, error) {
	eventType := EVENT_ROOM_CONNECTED
	message := RoomConnectedMessage{RoomName: room.Name, UserId: user.GetId(), EventName: eventType.GetName()}
	return &message, nil
}

func NewCheerAddedMessage(cheer cheers.Cheer) (*CheerAddedMessage, error) {
	eventType := EVENT_CHEER_ADDED
	message := CheerAddedMessage{Cheer: cheer, EventName: eventType.GetName()}
	return &message, nil
}
func NewRoomLastSecondsCheerCountMessage(count int) (*RoomLastSecondsCheerCountMessage, error) {
	eventType := EVENT_LAST_SECONDS_COUNT
	message := RoomLastSecondsCheerCountMessage{Count: count, EventName: eventType.GetName()}
	return &message, nil
}

func InitEventsSocket(conn *websocket.Conn, roomsServer *RoomsServer) {
	socket := eventsSocket{conn: conn, roomsServer: roomsServer,
		quitIntensityListener: make(chan struct{}),
		addedCheersChannel:    make(chan cheers.Cheer)}

	socket.processFirstMessage()
	socket.listenToClientMessages()
	socket.listenToRoomCheers()
	socket.sendRoomConnnectedMessage()
	socket.handleEventsAndSendMessages()
	socket.setCloseHandler()
}
