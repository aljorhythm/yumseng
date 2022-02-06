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
		log.Panicf("Room Registration Failed. Unrecognised Request: %s", string(msg))
	}

	roomName := joinRoomRequest.RoomName
	user, _ := s.roomsServer.UserService.GetUser(joinRoomRequest.UserId)
	s.user = user
	room := s.roomsServer.RoomServicer.GetOrCreateRoom(roomName)
	s.room = room
	clientId := fmt.Sprintf("user=%s uuid=%s", user.GetId(), uuid.New().String())
	s.clientId = clientId
	log.Printf("Room Registered: User %s in Room %s", s.user, s.room.Name)
	if err != nil {
		log.Panicf("error emitting room connected %s %#v", clientId, err)
	}
}

func (s *eventsSocket) listenToClientMessages() {
	go func() {
		log.Printf("[listenToClientMessages] EventsSocketId: %s", s.clientId)
		for {
			_, msg, err := s.conn.ReadMessage()

			if err != nil {
				log.Printf("[listenToClientMessages] Error reading socket message. EventsSocketId: %s Error: %#v", s.clientId, err)
				return
			}
			reader := bytes.NewReader(msg)
			newCheer := cheers.Cheer{}
			utils.DecodeJson(reader, &newCheer)

			log.Printf("EventsSocketId: %s Adding cheer %#v to %v", s.clientId, newCheer, s.room)
			s.roomsServer.RoomServicer.AddCheer(s.room, &newCheer, s.user)
		}
	}()
}

func (socket *eventsSocket) listenToRoomCheers() {
	log.Printf("EventsSocketId: %s Subscribing user %s to cheers in room %s ", socket.clientId, socket.user.GetId(), socket.room.Name)
	callback := func(args ...interface{}) {
		rawCheer := args[0]
		cheer, ok := rawCheer.(cheers.Cheer)
		if ok {
			log.Printf("EventsSocketId: %s Cheer Received from %s value: %s", socket.clientId, cheer.UserId, cheer.Value)
			socket.addedCheersChannel <- cheer
		} else {
			log.Panicf("EventsSocketId: %s Cheer Not Recognised %#v", socket.clientId, args)
		}
	}

	err := socket.roomsServer.AddCheerAddedListener(socket.room, socket.user, socket.clientId, callback)

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
	intensityTicker := time.NewTicker(250 * time.Millisecond)
	for {
		select {
		case cheer, more := <-socket.addedCheersChannel:
			if more {
				cheerAddedMessage, err := NewCheerAddedMessage(cheer)
				log.Printf("%s writing to socket cheer added from %s value: %s", socket.clientId, cheer.UserId, cheer.Value)
				err = socket.conn.WriteJSON(cheerAddedMessage)
				if err != nil {
					log.Panicf("client %s webSocket erroring write message %#v", socket.clientId, err)
				}
			} else {
				log.Printf("cheers channel is closed %s", socket.clientId)
			}
		case <-intensityTicker.C:
			intensity := socket.room.Intensity(func(cheer cheers.Cheer) bool {
				return cheer.Value == "yum"
			})
			message, _ := NewRoomIntensityMessage(intensity)
			err := socket.conn.WriteJSON(message)
			if err != nil {
				log.Printf("EventsSocketId: %s Error writing to connection %#v closing quit channel", socket.clientId, err)
				close(socket.quitIntensityListener)
			} else {
				//log.Printf("wrote to socket last seconds cheer intensity %s %f", socket.clientId, intensity)
			}
		case <-socket.quitIntensityListener:
			log.Printf("EventsSocketId: %s Stop sending cheer intensity", socket.clientId)
			intensityTicker.Stop()
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

func InitEventsSocket(conn *websocket.Conn, roomsServer *RoomsServer) {
	socket := eventsSocket{conn: conn, roomsServer: roomsServer,
		quitIntensityListener: make(chan struct{}),
		addedCheersChannel:    make(chan cheers.Cheer)}

	log.Printf("[InitEventsSocket] setCloseHandler")
	socket.setCloseHandler()
	log.Printf("[InitEventsSocket] processFirstMessage")
	socket.processFirstMessage()
	log.Printf("[InitEventsSocket] listenToClientMessages")
	socket.listenToClientMessages()
	log.Printf("[InitEventsSocket] listenToRoomCheers")
	socket.listenToRoomCheers()
	log.Printf("[InitEventsSocket] sendRoomConnectedMessage")
	socket.sendRoomConnnectedMessage()
	log.Printf("[InitEventsSocket] handleEventsAndSendMessages")
	socket.handleEventsAndSendMessages()
}
