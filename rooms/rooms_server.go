package rooms

import (
	"fmt"
	"github.com/aljorhythm/yumseng/objectstorage"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type RoomsServer struct {
	http.Handler
	RoomServicer
	UserService UserServicer
	RoomsServerOpts
}

type RoomsServerOpts struct {
	AllowOriginFunc func(r *http.Request) bool
}

func (roomsServer *RoomsServer) upgradeHttpToWs(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     roomsServer.AllowOriginFunc,
	}
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Panicf("webSocket upgrade error %#v", err)
	}
	return conn
}

func (roomsServer *RoomsServer) eventsWs(w http.ResponseWriter, r *http.Request) {
	conn := roomsServer.upgradeHttpToWs(w, r)
	InitEventsSocket(conn, roomsServer)
}

func NewRoomsServer(router *mux.Router, userService UserServicer, storage objectstorage.Storage, opts RoomsServerOpts) http.Handler {
	roomsServer := &RoomsServer{
		RoomServicer:    NewRoomsService(storage),
		UserService:     userService,
		RoomsServerOpts: opts,
	}

	router.Handle("/rooms", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "you are at rooms")
	}))

	router.Handle("/events", http.HandlerFunc(roomsServer.eventsWs))

	roomsServer.Handler = router
	return roomsServer
}
