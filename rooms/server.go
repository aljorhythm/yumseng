package rooms

import (
	"fmt"
	"github.com/aljorhythm/yumseng/cheers"
	"github.com/aljorhythm/yumseng/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"io/ioutil"
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

func (roomsServer *RoomsServer) roomUserImageHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[RoomsServer#roomUserImageHandler]")
	roomsService := roomsServer.RoomServicer
	userService := roomsServer.UserService

	vars := mux.Vars(r)
	roomId, _ := vars["room-id"]
	userId, _ := vars["user-id"]

	user, err := userService.GetUser(userId)

	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if r.Method == http.MethodGet {
		log.Printf("[RoomsServer#roomUserImageHandler] GET user-id %s room-id %s", userId, roomId)

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
		log.Printf("[RoomsServer#roomUserImageHandler] POST user-id %s room-id %s", userId, roomId)

		requestBytes, err := ioutil.ReadAll(r.Body)

		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		image, err := roomsService.AddCheerImage(r.Context(), roomId, user, requestBytes)

		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		responseBytes, err := utils.ToJson(image)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		fmt.Fprintf(w, string(responseBytes))
	} else {
		fmt.Fprintf(w, "operation does not exist")
	}
}

func (roomsServer *RoomsServer) roomUserCheersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if r.Method == http.MethodPost {
		roomId, _ := vars["room-id"]
		room := roomsServer.RoomServicer.GetOrCreateRoom(roomId)
		userId, _ := vars["user-id"]
		user, _ := roomsServer.UserService.GetUser(userId)
		cheer := &cheers.Cheer{}
		utils.DecodeJson(r.Body, cheer)

		roomsServer.RoomServicer.AddCheer(room, cheer, user)
		w.Write([]byte{})
	}
}

func NewRoomsServer(router *mux.Router, roomsService RoomServicer, userService UserServicer, opts RoomsServerOpts) http.Handler {
	roomsServer := &RoomsServer{
		RoomServicer:    roomsService,
		UserService:     userService,
		RoomsServerOpts: opts,
	}

	router.HandleFunc("/{room-id}/user/{user-id}/images",
		utils.ChainMiddlewares(
			roomsServer.roomUserImageHandler,
			utils.AddSetJsonHeaderMw),
	)

	router.HandleFunc("/{room-id}/user/{user-id}/cheers",
		utils.ChainMiddlewares(
			roomsServer.roomUserCheersHandler,
			utils.AddSetJsonHeaderMw),
	)

	router.Handle("/events", http.HandlerFunc(roomsServer.eventsWs))

	roomsServer.Handler = router
	return roomsServer
}
