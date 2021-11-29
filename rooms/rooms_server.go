package rooms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aljorhythm/yumseng/cheers"
	"github.com/aljorhythm/yumseng/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type RoomsServer struct {
	http.Handler
	RoomServicer
}

func getRoomFromRequest(r *http.Request) *Room {
	roomName := r.Header.Get("rooms-name")

	if roomName == "" {
		return NewRoom("default")
	} else {
		return NewRoom(roomName)
	}
}

func (roomsServer *RoomsServer) eventsWs(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Panicf("webSocket upgrade error %#v", err)
	}

	clientId := uuid.New().String()
	room := getRoomFromRequest(r)
	log.Printf("joined rooms %s | clientId : %s", room.Name, clientId)

	go func() {
		log.Printf("client %s listening for rooms %s cheers", clientId, room.Name)
		for {
			_, msg, err := conn.ReadMessage()

			if err != nil {
				log.Printf("error in connection room: %s client: %s err: %#v", room.Name, clientId, err)
				return
			}
			reader := bytes.NewReader(msg)
			newCheer := cheers.Cheer{}
			utils.DecodeJson(reader, &newCheer)
			log.Printf("adding cheer %#v", newCheer)
			roomsServer.RoomServicer.AddCheer(room, newCheer)
		}
	}()

	log.Printf("client %s listening", clientId)
	roomsServer.ListenCheer(room, clientId, func(args ...interface{}) {
		rawCheer := args[0]
		cheer, ok := rawCheer.(cheers.Cheer)

		if ok {
			log.Printf("%s writing to socket %#v", clientId, cheer)
			cheerBytes, err := json.Marshal(cheer)

			if err != nil {
				log.Panicf("client %s webSocket erroring decoding cheer %#v string: %s", clientId, err, string(cheerBytes))
				return
			}

			_, err = writeWs(conn, cheerBytes)

			if err != nil {
				log.Panicf("client %s webSocket erroring write message %#v", clientId, err)
			}
		} else {
			log.Panicf("cannot convert cheer %#v", args)
		}
	})

	conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("exiting listening client %s code %d text %s", clientId, code, text)
		roomsServer.StopListening(room, clientId)
		return nil
	})
}

func writeWs(conn *websocket.Conn, msg []byte) (int, error) {
	err := conn.WriteMessage(websocket.TextMessage, msg)

	if err != nil {
		return 0, err
	}

	return len(msg), nil
}

func NewRoomsServer(router *mux.Router) http.Handler {
	roomsServer := &RoomsServer{
		RoomServicer: NewRoomsService(),
	}

	router.Handle("/rooms", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "you are at rooms")
	}))

	router.Handle("/events", http.HandlerFunc(roomsServer.eventsWs))

	roomsServer.Handler = router
	return roomsServer
}
