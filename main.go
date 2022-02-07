package main

import (
	"embed"
	"fmt"
	"github.com/aljorhythm/yumseng/objectstorage"
	"github.com/aljorhythm/yumseng/ping"
	"github.com/aljorhythm/yumseng/rooms"
	"github.com/aljorhythm/yumseng/utils"
	"github.com/gorilla/mux"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rs/cors"
)

// web ui static assets
//go:embed webui/*
var webuiFs embed.FS

//go:embed react-ts-ui/build/*
var reactUiFs embed.FS

func generateUiHandler() (http.Handler, error) {
	uiFileSystem, err := fs.Sub(webuiFs, "webui")
	if err != nil {
		return nil, err
	}
	return http.FileServer(http.FS(uiFileSystem)), nil
}

func generateReactUiHandler() (http.Handler, error) {
	uiFileSystem, err := fs.Sub(reactUiFs, "react-ts-ui/build")
	if err != nil {
		return nil, err
	}
	return http.FileServer(http.FS(uiFileSystem)), nil
}

func generatePingHandler(tag string, timeArg time.Time) func(writer http.ResponseWriter, request *http.Request) {
	timeLocation, err := time.LoadLocation("Asia/Singapore")
	if err != nil {
		panic(err)
	}
	response := &ping.PingResponse{
		Tag:             tag,
		ServerStartTime: timeArg.In(timeLocation),
	}
	bytes, err := utils.ToJson(response)

	if err != nil {
		panic(err)
	}

	message := string(bytes)
	return func(writer http.ResponseWriter, request *http.Request) {
		_, err := fmt.Fprintf(writer, message)
		if err != nil {
		}
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

func allowSameOriginOrInList(list []string) func(*http.Request) bool {
	return func(r *http.Request) bool {
		return checkSameOrigin(r) || requestOriginIsInList(list, r)
	}
}

func main() {
	log.SetOutput(ioutil.Discard)

	router := mux.NewRouter()

	objectStorage := objectstorage.NewInmemoryStore()

	userService := DummyUserService{}
	roomsService := rooms.NewRoomsService(objectStorage)

	roomsSubrouter := router.PathPrefix("/rooms").Subrouter()
	opts := rooms.RoomsServerOpts{AllowOriginFunc: allowSameOriginOrInList(getAllowedOrigins())}
	rooms.NewRoomsServer(roomsSubrouter, roomsService, userService, opts)

	tag := getVersionTag()
	router.HandleFunc("/ping", generatePingHandler(tag, time.Now()))

	webuiHandler, err := generateUiHandler()
	if err != nil {
		log.Panicf("[main.go#main] Error generateUiHandler %s", err)
	}
	router.PathPrefix("/").Handler(webuiHandler)

	reactUiHandler, err := generateReactUiHandler()
	if err != nil {
		log.Panicf("[main.go#main] Error generateReactUiHandler %s", err)
	}
	router.PathPrefix("/archiveui").Handler(http.StripPrefix("/archiveui", reactUiHandler))

	port := getPort()
	portArg := fmt.Sprintf(":%s", port)
	log.Printf("Running router PORT=%s", portArg)

	httpCorsConfig := cors.New(cors.Options{
		AllowedOrigins:   getAllowedOrigins(),
		AllowCredentials: true,
	})
	log.Fatal(http.ListenAndServe(portArg, httpCorsConfig.Handler(router)))
}
