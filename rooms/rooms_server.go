package rooms

import (
	"fmt"
	"github.com/aljorhythm/yumseng/utils"
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

func writeError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "%#v", err)
}

func (roomsServer *RoomsServer) roomUserCheerHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[RoomsServer#roomUserCheerHandler]")
	if r.Method == http.MethodGet {
		vars := mux.Vars(r)
		roomId, _ := vars["room-id"]
		userId, _ := vars["user-id"]

		log.Printf("[RoomsServer#roomUserCheerHandler] user-id %s room-id %s", userId, roomId)

		roomsService := roomsServer.RoomServicer
		userService := roomsServer.UserService
		user, err := userService.GetUser(userId)

		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		images, err := roomsService.GetCheerImages(r.Context(), roomId, user)

		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		if bytes, err := utils.ToJson(images); err != nil {
			writeError(w, http.StatusBadRequest, err)
		} else {
			fmt.Fprintf(w, string(bytes))
		}
	} else if r.Method == http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "not implemented")
	} else {
		fmt.Fprintf(w, "operation does not exist")
	}
}

func NewRoomsServer(router *mux.Router, roomsService RoomServicer, userService UserServicer, opts RoomsServerOpts) http.Handler {
	roomsServer := &RoomsServer{
		RoomServicer:    roomsService,
		UserService:     userService,
		RoomsServerOpts: opts,
	}

	router.HandleFunc("/{room-id}/user/{user-id}/images", roomsServer.roomUserCheerHandler)

	router.Handle("/events", http.HandlerFunc(roomsServer.eventsWs))

	roomsServer.Handler = router
	return roomsServer
}
