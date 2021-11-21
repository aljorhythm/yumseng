package main

import (
	"embed"
	"errors"
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

func main() {
	webuiHandler, err := generateUiHandler()

	if err != nil {
		log.Panicf("[main.go#main] Error generateUiHandler %s", err)
	}

	http.Handle("/", webuiHandler)

	cheersService := cheers.NewService()
	http.HandleFunc("/cheers", setJsonResponseHeader(generateCheersHandler(cheersService)))

	tag := getVersionTag()
	http.HandleFunc("/ping", generatePingHandler(tag))

	port := getPort()
	portArg := fmt.Sprintf(":%s", port)
	log.Printf("Running server PORT=%s", portArg)
	log.Fatal(http.ListenAndServe(portArg, nil))
}
