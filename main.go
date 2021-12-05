package main

import (
	"embed"
	"errors"
	"fmt"
	"github.com/aljorhythm/yumseng/cheers"
	"github.com/aljorhythm/yumseng/objectstorage"
	"github.com/aljorhythm/yumseng/ping"
	"github.com/aljorhythm/yumseng/rooms"
	"github.com/aljorhythm/yumseng/utils"
	"github.com/gorilla/mux"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

// web ui static assets
//go:embed webui/*
var webuiFs embed.FS

// http errors
var (
	ERROR_UNHANDLED_HTTP_METHOD = errors.New("unhandled http method")
)

func generateUiHandler() (http.Handler, error) {
	uiFileSystem, err := fs.Sub(webuiFs, "webui")
	if err != nil {
		return nil, err
	}
	return http.FileServer(http.FS(uiFileSystem)), nil
}

func generateCheersHandler(service cheers.Servicer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			cheers := service.GetCheers()
			message, err := utils.ToJson(cheers)

			if err != nil {
				panic(err)
			}
			fmt.Fprintf(w, string(message))
		} else if r.Method == http.MethodPost {
			cheer := &cheers.Cheer{}
			err := utils.HttpRequestBodyToStruct(r, cheer)

			if err != nil {
				panic(err)
			}

			service.AddCheer(cheer)
			fmt.Fprintf(w, "{}")
		} else {
			panic(ERROR_UNHANDLED_HTTP_METHOD)
		}
	}
}

func generatePingHandler(tag string) func(writer http.ResponseWriter, request *http.Request) {
	response := &ping.PingResponse{
		Tag: tag,
	}
	bytes, error := utils.ToJson(response)

	if error != nil {
		panic(error)
	}

	message := string(bytes)
	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, message)
	}
}

func getPort() string {
	value := os.Getenv("PORT")
	if value == "" {
		log.Printf("Defaulting to default port 80")
		return "80"
	}
	return value
}

func getVersionTag() string {
	value := os.Getenv("TAG")
	if value == "" {
		log.Printf("Defaulting tag value to 'unknown-tag'")
		return "unknown-tag"
	}
	return value
}

func setJsonResponseHeader(handler func(writer http.ResponseWriter, request *http.Request)) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		handler(writer, request)
	}
}

type DummyUserService struct {
}

type DummyUser struct {
	id string
}

func (d DummyUser) GetId() string {
	return d.id
}

func (d DummyUserService) GetUser(id string) (rooms.User, error) {
	return DummyUser{id: id}, nil
}

var allowOriginFunc = func(r *http.Request) bool {
	log.Printf("There's a request! ")
	return true
}

func main() {
	router := mux.NewRouter()
	reactTsPort := 3000
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{fmt.Sprintf("http://localhost:%d", reactTsPort)},
		AllowCredentials: true,
	})
	objectStorage := objectstorage.NewInmemoryStore()

	cheersService := cheers.NewService()
	router.HandleFunc("/cheers", setJsonResponseHeader(generateCheersHandler(cheersService)))

	roomsSubrouter := router.PathPrefix("/rooms").Subrouter()

	userService := DummyUserService{}
	rooms.NewRoomsServer(roomsSubrouter, userService, objectStorage)

	tag := getVersionTag()
	router.HandleFunc("/ping", generatePingHandler(tag))

	webuiHandler, err := generateUiHandler()

	if err != nil {
		log.Panicf("[main.go#main] Error generateUiHandler %s", err)
	}

	// serve root html
	router.PathPrefix("/").Handler(webuiHandler)

	/* SocketIO */

	sioServer := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})

	sioServer.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())
		return nil
	})

	go func() {
		fmt.Println("socketio attempting to serve")
		if err := sioServer.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer sioServer.Close()

	http.Handle("/socket.io/", sioServer)

	port := getPort()
	portArg := fmt.Sprintf(":%s", port)
	log.Printf("Running router PORT=%s", portArg)
	log.Fatal(http.ListenAndServe(portArg, c.Handler(router)))
}
