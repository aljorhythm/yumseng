package rooms

import (
	"bytes"
	"fmt"
	"github.com/aljorhythm/yumseng/cheers"
	"github.com/aljorhythm/yumseng/utils"
	"github.com/aljorhythm/yumseng/utils/movingavg"
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
}

func getRoomFromRequest(r *http.Request) *Room {
	roomName := r.Header.Get("rooms-name")

	if roomName == "" {
		return NewRoom("default")
	} else {
		return NewRoom(roomName)
	}
}

type NowTime struct {
}

func (nowTime NowTime) Now() time.Time {
	return time.Now()
}

type CheerItem struct {
	cheers.Cheer
}

func (item CheerItem) GetTime() time.Time {
	return item.Cheer.ClientCreatedAt
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

	calculator := movingavg.NewCalculator(NowTime{})
	calculatorCheersChan := make(chan cheers.Cheer)

	log.Printf("client %s listening to room %s cheers", clientId, room.Name)
	roomsServer.ListenCheer(room, clientId, func(args ...interface{}) {
		rawCheer := args[0]
		cheer, ok := rawCheer.(cheers.Cheer)

		calculatorCheersChan <- cheer

		if ok {

			messageBytes, err := NewCheerAddedMessage(cheer)
			log.Printf("%s writing to socket %#v", clientId, cheer)

			if err != nil {
				log.Panicf("client %s webSocket erroring decoding cheer %#v string: %s", clientId, err, string(messageBytes))
				return
			}

			_, err = writeWs(conn, messageBytes)

			if err != nil {
				log.Panicf("client %s webSocket erroring write message %#v", clientId, err)
			}
		} else {
			log.Panicf("cannot convert cheer %#v", args)
		}
	})

	log.Printf("client %s listening to room %s cheer speed", clientId, room.Name)
	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case cheer, more := <-calculatorCheersChan:
				if more {
					log.Printf("adding cheer to calculator %#v", cheer)
					calculator.AddItem(CheerItem{cheer})
				} else {
					log.Printf("cheers channel is closed %s", clientId)
				}
			case <-ticker.C:
				count := calculator.CountFrom(time.Duration(1) * time.Second)
				message, err := NewRoomLastSecondsCheerCountMessage(count)

				if err != nil {
					log.Printf("err generating last seconds cheer count message %s %#v", clientId, err)
					continue
				} else {
					log.Printf("last seconds cheer count %s %d", clientId, count)
				}

				_, err = writeWs(conn, message)
				if err != nil {
					log.Printf("err writing to socket %#v closing quit channel %s", err, clientId)
					close(quit)
				} else {
					log.Printf("wrote to socket last seconds cheer count %s %d", clientId, count)
				}
			case <-quit:
				log.Printf("quit channel emitted stopping speed ticker %s", clientId)
				ticker.Stop()
				close(calculatorCheersChan)
				return
			}
		}
	}()

	conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("exiting listening client %s code %d text %s", clientId, code, text)
		roomsServer.StopListeningCheers(room, clientId)
		close(quit)
		close(calculatorCheersChan)
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
