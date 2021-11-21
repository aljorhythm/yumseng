package main

import (
	"embed"
	"fmt"
	"github.com/aljorhythm/yumseng/cheers"
	"github.com/aljorhythm/yumseng/ping"
	"github.com/aljorhythm/yumseng/utils"
	"io/fs"
	"log"
	"net/http"
	"os"
)

// web ui static assets
//go:embed webui/*
var webuiFs embed.FS

// tag for version / deployment control
//go:embed .tag
var tag string

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
		log.Printf("Resorting to default port 80")
		return "80"
	}
	return value
}

func main() {
	webuiHandler, err := generateUiHandler()

	if err != nil {
		log.Panicf("[main.go#main] Error generateUiHandler %s", err)
	}

	http.Handle("/", webuiHandler)

	cheersService := cheers.NewService()
	http.HandleFunc("/cheers", generateCheersHandler(cheersService))

	http.HandleFunc("/ping", generatePingHandler(tag))

	port := getPort()
	portArg := fmt.Sprintf(":%s", port)
	log.Printf("Running server PORT=%s", portArg)
	log.Fatal(http.ListenAndServe(portArg, nil))
}
